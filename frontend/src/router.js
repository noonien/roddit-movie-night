var app = require('ampersand-app');
var Router = require('ampersand-router');

import VotingPage from './pages/voting'


module.exports = Router.extend({
    routes: {
        '': 'voting',
        '(*path)': 'catchAll'
    },

    // ------- ROUTE HANDLERS ---------
    voting() {
        app.trigger('page', new VotingPage({
            model: app.me,
            collection: app.movies
        }));
    },
    catchAll() {
        this.redirectTo('');
    }
});
