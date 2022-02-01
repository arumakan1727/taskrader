const API_ORIGIN = location.origin;
const autoFetch = true;

const DAY_TABLE = ['日', '月', '火', '水', '木', '金', '土'];

window.onbeforeunload = () => {
  fetch(API_ORIGIN + '/ses/leave', { method: 'POST' });
}

function zeroPadding(value, width) {
  return String(value).padStart(width, '0');
}

Vue.component('login-status', {
  props: {
    logined: Boolean,
    verifying: Boolean,
  },
  template: `

  <div v-if="verifying">
    <md-progress-spinner style="vertical-align:middle" :md-diameter="24" :md-stroke="3" md-mode="indeterminate"></md-progress-spinner>
    <span class="ml-1">ログイン試行中...</span>
  </div>
  <div v-else-if="logined" class="accepted">
    <md-icon>check_circle</md-icon> <span>登録済み</span>
  </div>
  <div v-else class="not-accepted">
    <md-icon>warning</md-icon> <span>未登録</span>
  </div>
  `
});

Vue.use(VueMaterial.default);
const app = new Vue({
  el: '#app',
  data: () => ({
    state: {
      initialized: false,
      menuVisible: false,
      tab: 'HOME',
      fetchingAss: false,
      puttingGakujoAuth: false,
      puttingEdstemAuth: false,
      puttingTeamsAuth: false,
    },
    sorting: {
      key: 'due',
      ord: 'asc',
    },
    assignments: [],
    assErrors: [],
    auth: {
      gakujo: {
        username: "",
        password: "",
        accepted: false,
        errmsg: "",
      },
      edstem: {
        email: "",
        password: "",
        accepted: false,
        errmsg: "",
      },
      teams: {
        email: "",
        password: "",
        accepted: false,
        errmsg: "",
      },
    }
  }),
  mounted() {
    this.fetchAuth().then(() => {
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
    gakujoAuthEmpty() {
      return !this.auth.gakujo.username || !this.auth.gakujo.password;
    },
    edstemAuthEmpty() {
      return !this.auth.edstem.email || !this.auth.edstem.password;
    },
    teamsAuthEmpty() {
      return !this.auth.teams.email || !this.auth.teams.password;
    },
    nothingLogined() {
      const a = this.auth;
      return !a.gakujo.accepted && !a.edstem.accepted && !a.teams.accepted;
    },
  },
  methods: {
    toggleMenu() {
      this.state.menuVisible = !this.state.menuVisible;
    },
    switchTo(tabID) {
      this.state.tab = tabID;
    },
    fetchAss() {
      if (this.state.fetchingAss) return;
      this.state.fetchingAss = true;
      return axios.get(API_ORIGIN + '/api/assignments')
        .then(resp => {
          const ass = resp.data.assignments;
          const now = new Date();
          ass.forEach(a => {
            a.due = new Date(a.due);
            a.fmtDue = this.fmtDue(a.due, now);
            a.remHours = (a.due - now) / 36e5;
          });
          this.assignments = ass;
          this.assErrors = resp.data.errors;
        })
        .catch(err => {
          console.error(err);
        })
        .finally(() => {
          this.state.fetchingAss = false;
        })
    },
    fetchAuth() {
      return axios.get(API_ORIGIN + '/api/auth')
        .then(resp => {
          const g = resp.data.gakujo;
          const e = resp.data.edstem;
          const t = resp.data.teams;
          this.auth.gakujo.username = g.username;
          this.auth.gakujo.password = g.password;
          this.auth.edstem.email = e.email;
          this.auth.edstem.password = e.password;
          this.auth.teams.email = t.email;
          this.auth.teams.password = t.password;
          this.auth.gakujo.accepted = Boolean(this.auth.gakujo.username && this.auth.gakujo.password);
          this.auth.edstem.accepted = Boolean(this.auth.edstem.email && this.auth.edstem.password);
          this.auth.teams.accepted = Boolean(this.auth.teams.email && this.auth.teams.password);
        })
        .catch(err => {
          console.error(err);
        })
    },
    putGakujoAuth() {
      if (this.gakujoAuthEmpty) {
        this.auth.gakujo.errmsg = "Username または Password が空です";
        return;
      }
      this.state.puttingGakujoAuth = true;
      this.auth.gakujo.errmsg = "";

      axios.put(API_ORIGIN + '/api/auth/gakujo', {
        username: this.auth.gakujo.username,
        password: this.auth.gakujo.password,
      })
        .then(resp => {
          if (resp.data.errmsg) {
            this.auth.gakujo.errmsg = resp.data.errmsg;
          } else {
            this.auth.gakujo.accepted = true;
            if (this.assignments.length == 0) this.fetchAss();
          }
        })
        .finally(() => {
          this.state.puttingGakujoAuth = false;
        })
    },
    putEdstemAuth() {
      if (this.edstemAuthEmpty) {
        this.auth.edstem.errmsg = "Email または Password が空です";
        return;
      }
      this.state.puttingEdstemAuth = true;
      this.auth.edstem.errmsg = "";

      axios.put(API_ORIGIN + '/api/auth/edstem', {
        email: this.auth.edstem.email,
        password: this.auth.edstem.password,
      })
        .then(resp => {
          if (resp.data.errmsg) {
            this.auth.edstem.errmsg = resp.data.errmsg;
          } else {
            this.auth.edstem.accepted = true;
            if (this.assignments.length == 0) this.fetchAss();
          }
        })
        .finally(() => {
          this.state.puttingEdstemAuth = false;
        })
    },
    putTeamsAuth() {
      if (this.teamsAuthEmpty) {
        this.auth.teams.errmsg = "Email または Password が空です";
        return;
      }
      this.state.puttingTeamsAuth = true;
      this.auth.teams.errmsg = "";

      axios.put(API_ORIGIN + '/api/auth/teams', {
        email: this.auth.teams.email,
        password: this.auth.teams.password,
      })
        .then(resp => {
          if (resp.data.errmsg) {
            this.auth.teams.errmsg = resp.data.errmsg;
          } else {
            this.auth.teams.accepted = true;
            if (this.assignments.length == 0) this.fetchAss();
          }
        })
        .finally(() => {
          this.state.puttingTeamsAuth = false;
        })
    },
    fmtDue(d, now) {
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

      let res = `${MM + 1}月${DD}日 (${day}) ${zeroPadding(hh, 2)}:${zeroPadding(mm, 2)}`;
      return (YYYY !== now.getFullYear()) ? (YYYY + '年' + res) : res;
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
});
