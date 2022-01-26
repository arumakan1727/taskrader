package cred

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/denisbrodbeck/machineid"
)

type Credential struct {
	Gakujo Gakujo `json:"gakujo"`
	EdStem EdStem `json:"edstem"`
	Teams  Teams  `json:"teams"`
}

type Gakujo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type EdStem struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Teams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// filepath から暗号化されたデータを読み出して復号して credential を生成する。
// データのフォーマットに異常があった場合はファイルを削除する。
func LoadFromFile(filepath string) (*Credential, error) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	bs, err = decrypt(bs)
	if err != nil {
		os.Remove(filepath)
		return nil, err
	}

	var c Credential
	if err := json.Unmarshal(bs, &c); err != nil {
		os.Remove(filepath)
		return nil, err
	}
	return &c, nil
}

// filepath からの読み込みを試行し、成功すればその結果をそのまま返し、失敗したら空の credential を返す。
func LoadFromFileOrEmpty(filepath string) *Credential {
	c, err := LoadFromFile(filepath)
	if err != nil {
		return &Credential{}
	}
	return c
}

// filepath の親ディレクトリを再帰的に作成してから filepath に書き出す。
// ディレクトリ作成時のエラーは無視される。
func (c *Credential) SaveToFile(filepath string) error {
	bs, err := json.Marshal(c)
	if err != nil {
		return err
	}

	bs, err = encrypt(bs)
	if err != nil {
		return err
	}

	os.MkdirAll(path.Dir(filepath), 0700)

	return ioutil.WriteFile(filepath, bs, 0600)
}


func encrypt(plaintext []byte) ([]byte, error) {
	key := getKey()
	nonce := getNonce()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(ciphertext []byte) ([]byte, error) {
	key := getKey()
	nonce := getNonce()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}

func getKey() []byte {
	var machineID []byte

	if id, err := machineid.ID(); err == nil {
		machineID = []byte(id)
	} else {
		dir, _ := os.UserHomeDir()
		slug := make([]byte, 0, 40)
		slug = append(slug, []byte(dir)...)

		var x, a, b, c byte = 117, 5, 3, 6
		for len(slug) < 40 {
			x ^= x >> a
			x ^= x << b
			x ^= x >> c
			slug = append(slug, x)
		}
		machineID = slug
	}

	mac := hmac.New(sha256.New, machineID)
	mac.Write([]byte("taskrader"))
	return mac.Sum(nil)
}

func getNonce() []byte {
	return []byte{47, 239, 81, 211, 110, 210, 75, 254, 227, 145, 206, 210}
}

// 環境変数から Credential を生成する。
// .env は自動では読み込まないので適宜 godotenv.Load() すること。
// 該当する名前の環境変数が存在しない場合は、そのメンバ変数の値は空文字列となる。
func LoadFromEnv() *Credential {
	return &Credential{
		Gakujo: Gakujo{
			Username: os.Getenv("GAKUJO_USERNAME"),
			Password: os.Getenv("GAKUJO_PASSWORD"),
		},
		EdStem: EdStem{
			Email:    os.Getenv("EDSTEM_EMAIL"),
			Password: os.Getenv("EDSTEM_PASSWORD"),
		},
		Teams: Teams{
			Email:    os.Getenv("TEAMS_EMAIL"),
			Password: os.Getenv("TEAMS_PASSWORD"),
		},
	}
}

type ErrEmpty struct {
	FieldName string
}

func (e *ErrEmpty) Error() string {
	return fmt.Sprintf("%s is empty", e.FieldName)
}

func newErrEmpty(fieldName string) ErrEmpty {
	return ErrEmpty{
		FieldName: fieldName,
	}
}

func (c *Credential) CheckEmptyField() []ErrEmpty {
	errs := make([]ErrEmpty, 0)
	if c.Gakujo.Username == "" {
		errs = append(errs, newErrEmpty("Gakujo.Username"))
	}
	if c.Gakujo.Password == "" {
		errs = append(errs, newErrEmpty("Gakujo.Username"))
	}
	if c.EdStem.Email == "" {
		errs = append(errs, newErrEmpty("EdStem.Email"))
	}
	if c.EdStem.Password == "" {
		errs = append(errs, newErrEmpty("EdStem.Password"))
	}
	if c.Teams.Email == "" {
		errs = append(errs, newErrEmpty("Teams.Email"))
	}
	if c.Teams.Password == "" {
		errs = append(errs, newErrEmpty("Teams.Password"))
	}
	return errs
}

func (c *Credential) AbortIfEmptyFieldExists() {
	if errs := c.CheckEmptyField(); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, "Error: "+e.Error())
		}
		os.Exit(1)
	}
}
