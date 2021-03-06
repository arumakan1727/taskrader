package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arumakan1727/taskrader/assets"
	"github.com/arumakan1727/taskrader/pkg/assignment"
	"github.com/arumakan1727/taskrader/pkg/clients/edstem"
	"github.com/arumakan1727/taskrader/pkg/clients/gakujo"
	"github.com/arumakan1727/taskrader/pkg/clients/teams"
	"github.com/arumakan1727/taskrader/pkg/config"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/arumakan1727/taskrader/pkg/view"
	"github.com/gin-gonic/gin"
)

type AssignmentsSupplyer = func(auth *cred.Credential) ([]*assignment.Assignment, []*assignment.Error)

var (
	authFilepath    string
	errAuthFilepath error
)

func init() {
	authFilepath, errAuthFilepath = config.TaskraderCredentialPath()
}

func NewEngine(assignmentsSupplyer AssignmentsSupplyer) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/taskrader")
	})

	apiRouter := r.Group("/api")
	{
		apiRouter.GET("/assignments", funcGetAssignments(assignmentsSupplyer))
		apiRouter.GET("/auth/status", getAuthStatus)
		apiRouter.GET("/auth", getAuth)
		apiRouter.PUT("/auth/gakujo", putAuthGakujo)
		apiRouter.PUT("/auth/edstem", putAuthEdstem)
		apiRouter.PUT("/auth/teams", putAuthTeams)
	}

	sesRouter := r.Group("/ses")
	{
		entered := make(chan bool, 1)

		sesRouter.POST("/enter", func(c *gin.Context) {
			c.String(200, "")
			entered <- true
		})

		sesRouter.POST("/leave", func(c *gin.Context) {
			// clear the channel
			for len(entered) > 0 {
				<-entered
			}

			select {
			case <-entered:
				c.String(200, "Detected reloading")

			case <-time.After(2 * time.Second):
				fmt.Println("Shutdown the taskrader server.")
				os.Exit(0)
			}
		})

	}

	return r
}

func AddAssetsRoute(r *gin.Engine) {
	r.GET("/taskrader", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", assets.IndexHTML())
	})
	r.GET("/file/main.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/javascript", assets.MainJS())
	})
}

func funcGetAssignments(assignmentsSupplyer AssignmentsSupplyer) func(*gin.Context) {
	return func(c *gin.Context) {
		if errAuthFilepath != nil {
			respAuthPathErr(c)
			return
		}

		auth := cred.LoadFromFileOrEmpty(authFilepath)

		ass, errs := assignmentsSupplyer(auth)
		view.SortAssignments(ass)
		resp := RespAssignmentsAndErrors{
			Ass:    ass,
			Errors: make([]RespAssErr, 0, len(errs)),
		}
		for _, e := range errs {
			resp.Errors = append(resp.Errors, RespAssErr{
				Origin:  string(e.Origin),
				Message: e.Err.Error(),
			})
		}
		c.JSON(http.StatusOK, &resp)
	}
}

func getAuth(c *gin.Context) {
	if errAuthFilepath != nil {
		respAuthPathErr(c)
		return
	}
	auth := cred.LoadFromFileOrEmpty(authFilepath)
	c.JSON(http.StatusOK, auth)
}

func getAuthStatus(c *gin.Context) {
	if errAuthFilepath != nil {
		respAuthPathErr(c)
		return
	}
	auth := cred.LoadFromFileOrEmpty(authFilepath)
	resp := RespLoginStatus{
		GakujoLogined: !auth.Gakujo.IsEmpty(),
		EdstemLogined: !auth.EdStem.IsEmpty(),
		TeamsLogined:  !auth.Teams.IsEmpty(),
	}
	c.JSON(http.StatusOK, resp)
}

func putAuthGakujo(c *gin.Context) {
	if errAuthFilepath != nil {
		respAuthPathErr(c)
		return
	}
	var gakujoCred cred.Gakujo
	if err := c.BindJSON(&gakujoCred); err != nil {
		return
	}

	err := gakujo.NewClient().Login(gakujoCred.Username, gakujoCred.Password)
	if err != nil {
		respLoginErr(err, c)
		return
	}

	auth := cred.LoadFromFileOrEmpty(authFilepath)
	auth.Gakujo = gakujoCred
	err = auth.SaveToFile(authFilepath)
	resp500SimpleErrOr200EmptyErr(err, c)
}

func putAuthEdstem(c *gin.Context) {
	if errAuthFilepath != nil {
		respAuthPathErr(c)
		return
	}
	var edstemCred cred.EdStem
	if err := c.BindJSON(&edstemCred); err != nil {
		return
	}

	err := edstem.NewClient().Login(edstemCred.Email, edstemCred.Password)
	if err != nil {
		respLoginErr(err, c)
		return
	}

	auth := cred.LoadFromFileOrEmpty(authFilepath)
	auth.EdStem = edstemCred
	err = auth.SaveToFile(authFilepath)
	resp500SimpleErrOr200EmptyErr(err, c)
}

func putAuthTeams(c *gin.Context) {
	if errAuthFilepath != nil {
		respAuthPathErr(c)
		return
	}
	var teamsCred cred.Teams
	if err := c.BindJSON(&teamsCred); err != nil {
		return
	}

	err := teams.Login(teamsCred.Email, teamsCred.Password, log.New(io.Discard, "", 0))
	if err != nil {
		respLoginErr(err, c)
		return
	}

	auth := cred.LoadFromFileOrEmpty(authFilepath)
	auth.Teams = teamsCred
	err = auth.SaveToFile(authFilepath)
	resp500SimpleErrOr200EmptyErr(err, c)
}

func respAuthPathErr(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, RespSimpleErr{Message: errAuthFilepath.Error()})
}

func resp500SimpleErrOr200EmptyErr(err error, c *gin.Context) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespSimpleErr{Message: err.Error()})
	} else {
		c.JSON(http.StatusOK, RespSimpleErr{Message: ""})
	}
}

func respLoginErr(err error, c *gin.Context) {
	if err == nil {
		c.JSON(200, RespSimpleErr{})
		return
	}

	resp := RespSimpleErr{}
	switch err := err.(type) {
	case *gakujo.ErrUsernameOrPasswdWrong:
		resp.Message = fmt.Sprintf("??????????????????: ???????????????????????????????????????????????????????????????????????? (username: '%s')", err.Username)
	case *edstem.ErrEmailOrPasswdWrong:
		resp.Message = fmt.Sprintf("??????????????????: ????????????????????????????????????????????????????????????????????????????????? (email: '%s')", err.Email)
	case *teams.ErrEmailOrPasswdWrong:
		resp.Message = fmt.Sprintf("??????????????????: ????????????????????????????????????????????????????????????????????????????????? (email: '%s')", err.Email)
	case *teams.ErrAlreadyLogined:
		resp.Message = "??????????????????????????????"
	default:
		resp.Message = err.Error()
	}

	c.JSON(200, &resp)
}
