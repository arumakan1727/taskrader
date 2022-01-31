const API_ORIGIN = location.origin;
const autoFetch = true;

const DAY_TABLE =  ['日', '月', '火', '水', '木', '金', '土'];

function zeroPadding(value, width) {
  return String(value).padStart(width, '0');
}

Vue.use(VueMaterial.default);
const app = new Vue({
  el: '#app',
  mounted() {
    this.fetchLoginStatus().then(() => {
      if (!this.nothingLogined && autoFetch) {
        this.fetchAss().then(() => {
          this.state.initialized = true;
        });
      } else {
        this.state.initialized = true;
      }
    })
  },
  computed: {
    nothingLogined() {
      return !this.logined.gakujo && !this.logined.edstem && !this.logined.teams;
    }
  },
  methods: {
    toggleMenu() {
      this.state.menuVisible = !this.state.menuVisible;
    },
    switchTo(tabID) {
      this.state.tab = tabID;
    },
    fetchAss() {
      this.state.fetchingAss = true;
      return axios.get(API_ORIGIN + '/api/assignments')
        .then(resp => {
          let ass = resp.data.assignments;
          ass.forEach(a => a.due = new Date(a.due));
          this.assignments = ass;
          this.assErrors = resp.data.errors;
        })
        .catch(err => {
          console.error(err);
        })
        .finally(() => {
          this.state.fetchingAss = false;
          this.state.now = new Date();
        })
    },
    fetchLoginStatus() {
      return axios.get(API_ORIGIN + '/api/auth/status')
        .then(resp => {
          this.logined = resp.data;
        })
        .catch(err => {
          console.error(err);
        })
    },
    fmtDue(d) {
      let YYYY = d.getFullYear();
      if (YYYY >= 9999) return "締切不明"

      let MM = d.getMonth();
      let DD = d.getDate();
      let hh = d.getHours();
      let mm = d.getMinutes();
      let day = DAY_TABLE[d.getDay()];

      if (hh === 0) {
        hh = 24;
        const yesterday = new Date(YYYY, MM, DD - 1);
        YYYY = yesterday.getFullYear();
        MM = yesterday.getMonth();
        DD = yesterday.getDate();
      }

      let res = `${MM+1}月${DD}日 (${day}) ${zeroPadding(hh, 2)}:${zeroPadding(mm, 2)}`;
      return (YYYY !== this.state.now.getFullYear()) ? (YYYY + '年' + res) : res;
    },
    sortAssignments(ass) {
      const key = this.sorting.key;
      const k = (this.sorting.ord === 'desc') ? -1 : 1;
      if (key === 'due') {
        return ass.sort((a, b) => k * (a[key] - b[key]))
      } else {
        return ass.sort((a, b) => k * a[key].localeCompare(b[key]))
      }
    },
  },
  data: () => ({
    state: {
      initialized: false,
      menuVisible: false,
      tab: 'HOME',
      fetchingAss: false,
      puttingAuth: false,
      now: new Date(),
    },
    sorting: {
      key: 'due',
      ord: 'asc',
    },
    assignments: [
      {
        origin: "学情",
        title: "最終レポート",
        course: "計算理論",
        due: '2022-01-30T23:55:00+0900',
      },
      {
        origin: "EdStem",
        title: "小レポート99",
        course: "データベースシステム論",
        due: '9999-12-31T23:59:59Z',
      },
    ],
    assErrors: [
      {
        origin: "Teams",
        message: "email, password が空です"
      },
    ],
    logined: {
      gakujo: false,
      edstem: false,
      teams: false,
    }
  }),
});
