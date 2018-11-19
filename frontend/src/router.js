import app from 'ampersand-app'
import Router from 'ampersand-router'

import Movies from './models/movies'
import Poll from './models/poll'

import VotingPage from './pages/voting'


module.exports = Router.extend({
    routes: {
        '': 'voting',
        '(*path)': 'catchAll'
    },

    // ------- ROUTE HANDLERS ---------
    voting() {
        let poll = new Poll()

        poll.fetch({
            success(poll) {
                for (let vote of poll.votes) {
                    let movie = poll.movies.findWhere({ id: vote })
                    movie.selected = true
                }

                app.trigger('page', new VotingPage({
                    model: poll,
                }));
            }
        })
        
    },
    catchAll() {
        this.redirectTo('');
    }
});
