import PageView from './base'
import $ from 'jquery'
import _ from 'lodash'
import toastr from 'toastr'
import showdown from 'showdown'
let markdown = new showdown.Converter()

import votingTemplate from '../../templates/pages/voting'

import MovieView from '../views/movie'

export default PageView.extend({
    pageTitle: 'Movie Night',
    template: votingTemplate,
    bindings: {
        'model.poll.name': {
            type: 'text',
            hook: 'name',
        },
        'model.poll.info': {
            type: (el, val) => {
                el.innerHTML = markdown.makeHtml(val)
            },
            hook: 'info',
        },
    },
    events: {
        'click [data-hook~=action-vote]': 'vote',
        'click [data-hook~=action-suggest-movie]': 'suggestMovie',
    },
    initialize() {
        setInterval(() => {
            this.fetchModel()
        }, 15000)
    },
    render: function () {
        this.renderWithTemplate()
        this.renderCollection(this.model.movies,
            MovieView,
            this.queryByHook('movie-list'))

        if (!this.model.poll.info) {
            $(this.queryByHook('info-container')).hide()
        }
        
        this.fetchModel()
    },
    vote(id, voted) {
        let votes = this.model.movies.where({ selected: true })
            .map(movie => movie.id)

        // increment vote locally
        let votedMovie = this.model.movies.findWhere({ id })
        let increment = voted ? 1 : -1
        votedMovie.votes += increment

        $.ajax({
            url: '/api/polls/latest/vote',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                votes
            })
        })
        .fail(err => {
            toastr.error('S-a produs o eroare, incearca din nou. Sorry!')
        })
    },
    suggestMovie() {
        let $input = $(this.queryByHook('suggest-input'))
        let url = $input.val()

        $.ajax({
            url: '/api/polls/latest/suggest',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                url
            })
        })
        .then(() => {
            $input.val(null)
            toastr.success('Ai propus cu succes. Merci!')
        })
        .fail(err => {
            toastr.error('S-a produs o eroare, incearca din nou. Sorry!')
        })
        .always(() => {
            this.fetchModel();
        })
    },
    fetchModel() {
        this.model.fetch({ 
            success: () => {
                for (let vote of this.model.votes) {
                    let movie = this.model.movies.findWhere({ id: vote })
                    movie.selected = true
                }
            }
        });
        return false;
    },
});
