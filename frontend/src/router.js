var app = require('ampersand-app');
var Router = require('ampersand-router');

import VotingPage from './pages/voting'


module.exports = Router.extend({
    routes: {
        '': 'voting',
        '/results': 'results',
        '(*path)': 'catchAll'
    },

    // ------- ROUTE HANDLERS ---------
    voting() {
        app.trigger('page', new VotingPage({
            model: app.me,
            collection: app.movies
        }));
    },
    results() {
        app.trigger('page', new VotingPage({
            model: app.me,
            collection: app.movies,
            readonly: true,
        }));
    },
    catchAll() {
        this.redirectTo('');
    }
});
