import Vuex from 'vuex';
import { mount, createLocalVue } from '@vue/test-utils';
import Vuetify from 'vuetify';
import DeviceWarning from '@/components/device/DeviceWarning';

describe('DeviceWarning', () => {
  const localVue = createLocalVue();
  const vuetify = new Vuetify();
  localVue.use(Vuex);

  let wrapper;

  const dialog = false;
  const deviceWarning = false;
  const devicesSelected = [];

  const hostname = 'localhost';
  const url = `http://${hostname}/settings/billing`;

  const filter = [];

  const items = [
    {
      title: 'Suggested Devices',
      action: 'suggestedDevices',
    },
    {
      title: 'All devices',
      action: 'allDevices',
    },
  ];

  const devices = [
    {
      uid: 'a582b47a42d',
      name: '39-5e-2a',
      identity: {
        mac: '00:00:00:00:00:00',
      },
      info: {
        id: 'linuxmint',
        pretty_name: 'Linux Mint 19.3',
        version: '',
      },
      public_key: '----- PUBLIC KEY -----',
      tenant_id: '00000000',
      last_seen: '2020-05-20T18:58:53.276Z',
      online: false,
      namespace: 'user',
      status: 'accepted',
    },
    {
      uid: 'a582b47a42e',
      name: '39-5e-2b',
      identity: {
        mac: '00:00:00:00:00:00',
      },
      info: {
        id: 'linuxmint',
        pretty_name: 'Linux Mint 19.3',
        version: '',
      },
      public_key: '----- PUBLIC KEY -----',
      tenant_id: '00000001',
      last_seen: '2020-05-20T19:58:53.276Z',
      online: true,
      namespace: 'user',
      status: 'accepted',
    },
  ];

  const store = new Vuex.Store({
    namespaced: true,
    state: {
      deviceWarning,
      devicesSelected,
      filter,
      devices,
    },
    getters: {
      'devices/getDeviceWarning': (state) => state.deviceWarning,
      'devices/getDevicesSelected': (state) => state.devicesSelected,
      'devices/getFilter': (state) => state.filter,
      'devices/list': (state) => state.devices,
    },
    actions: {
      'devices/getDevicesMostUsed': () => {},
      'devices/fetch': () => {},
      'snackbar/showSnackbarDeviceChoice': () => {},
      'devices/postDevicesChoice': () => {},
      'stats/get': () => {},
      'devices/setDeviceWarning': () => {},
      'snackbar/showSnackbarErrorAssociation': () => {},
      'snackbar/showSnackbarErrorLoading': () => {},
    },
  });

  describe('Dialog is closes', () => {
    beforeEach(() => {
      wrapper = mount(DeviceWarning, {
        store,
        localVue,
        stubs: ['fragment'],
        vuetify,
      });
    });

    ///////
    // Component Rendering
    //////

    it('Is a Vue instance', () => {
      document.body.setAttribute('data-app', true);
      expect(wrapper).toBeTruthy();
    });
    it('Renders the component', () => {
      expect(wrapper.html()).toMatchSnapshot();
    });

    ///////
    // Data and Props checking
    //////

    it('Compare data with default value', () => {
      expect(wrapper.vm.hostname).toEqual(hostname);
      expect(wrapper.vm.action).toEqual(items[0].action);
      expect(wrapper.vm.items).toEqual(items);
      expect(wrapper.vm.dialog).toEqual(dialog);
    });
    it('Process data in the computed', () => {
      expect(wrapper.vm.show).toEqual(false);
      expect(wrapper.vm.showTooltip).toEqual(false);
      expect(wrapper.vm.equalThreeDevices).toEqual(false);
    });
    it('Process data in methods', () => {
      expect(wrapper.vm.url()).toEqual(url);

      wrapper.vm.close();
      expect(wrapper.vm.show).toEqual(false);
    });

    //////
    // HTML validation
    //////

    it('Renders the template with components', async () => {
      expect(wrapper.find('[data-test="deviceListChoice-component"]').exists()).toEqual(false);
    });
    it('Renders the template with data', () => {
      expect(wrapper.find('[data-test="deviceWarning-dialog"]').exists()).toEqual(false);
      expect(wrapper.find('[data-test="close-btn"]').exists()).toEqual(false);
      expect(wrapper.find('[data-test="accept-btn"]').exists()).toEqual(false);
    });
  });
});
