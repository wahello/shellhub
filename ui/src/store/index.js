import Vue from 'vue';
import Vuex from 'vuex';
import stats from '@/store/modules/stats';
import sessions from '@/store/modules/sessions';
import auth from '@/store/modules/auth';
import devices from '@/store/modules/devices';
import modals from '@/store/modules/modals';
import snackbar from '@/store/modules/snackbar';
import firewallrules from '@/store/modules/firewall_rules';
import publickeys from '@/store/modules/public_keys';
import privatekeys from '@/store/modules/private_keys';
import notifications from '@/store/modules/notifications';
import users from '@/store/modules/users';
import security from '@/store/modules/security';
import namespaces from '@/store/modules/namespaces';
import boxs from '@/store/modules/boxs';
import mobile from '@/store/modules/mobile';
import tags from '@/store/modules/tags';
import spinner from '@/store/modules/spinner';
import billing from '@/store/modules/billing';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    devices,
    modals,
    snackbar,
    stats,
    sessions,
    auth,
    firewallrules,
    publickeys,
    privatekeys,
    notifications,
    users,
    security,
    namespaces,
    boxs,
    mobile,
    tags,
    spinner,
    billing,
  },
});
