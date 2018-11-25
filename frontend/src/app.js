
import app from 'ampersand-app'
import _ from 'lodash'
import $ from 'jquery'
import Router from './router'
import MainView from './views/main'
import toastr from 'toastr'
import * as Fingerprint2 from 'fingerprintjs2'
import { getCookie } from './utils'

import '../stylesheets/app.styl'
import '../stylesheets/vendor/bootstrap-darkly.min.css'
import 'pretty-checkbox/dist/pretty-checkbox'
import 'toastr/build/toastr.css'

let FP2 = Fingerprint2.default()

let version = '14:41'

toastr.options = {
    "positionClass": "toast-bottom-right",
}

// attach our app to `window` so we can
// easily access it from the console.
window.app = app
window.$ = $
window._ = _

// Extends our main app singleton
app.extend({
    router: new Router(),
    // This is where it all starts
    init() {
        // Create and attach our main view
        this.mainView = new MainView({
            el: document.body
        });

        // this kicks off our backbutton tracking (browser history)
        // and will cause the first matching handler in the router
        // to fire.
        this.router.history.start({ pushState: true });
    },
    // This is a helper for navigating around the app.
    // this gets called by a global click handler that handles
    // all the <a> tags in the app.
    // it expects a url pathname for example: "/costello/settings"
    navigate(page) {
        var url = (page.charAt(0) === '/') ? page.slice(1) : page;
        this.router.history.navigate(url, {
            trigger: true
        });
    }
});

if (!getCookie('fp')) {
    const setFP = () => {
        FP2.get(fp => {
            console.log('SET FP')
            document.cookie = `fp=${fp}` 
            app.init()
        })
    }

    if (window.requestIdleCallback) {
        requestIdleCallback(setFP)
    } else {
        setTimeout(setFP, 500)
    }
} else {
    console.log('INIT')
    app.init()
}
