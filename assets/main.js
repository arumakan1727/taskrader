const API_ORIGIN = location.origin;
const autoFetch = false;

Vue.use(VueMaterial.default);
const App = new Vue({
  el: '#app',
  mounted() {
    this.fetchLoginStatus().then(() => {
      this.state.initialized = true;
      if (!this.nothingLogined && autoFetch) this.fetchAss();
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
      axios.get(API_ORIGIN+'/api/assignments')
        .then(resp => {
          this.assignments = resp.data.assignments;
          this.assErrors = resp.data.errors;
        })
        .catch(err => {
          console.error(err);
        })
        .finally(() => {
          this.state.fetchingAss = false;
        })
    },
    fetchLoginStatus() {
      return axios.get(API_ORIGIN+'/api/auth/status')
        .then(resp => {
          console.log(resp.data);
          this.logined = resp.data;
        })
        .catch(err => {
          console.error(err);
        })
    },
  },
  data: () => ({
    state: {
      initialized: false,
      menuVisible: false,
      tab: 'HOME',
      fetchingAss: false,
      puttingAuth: false,
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
