import PageView from './base'
import $ from 'jquery'
import toastr from 'toastr'

import votingTemplate from '../../templates/pages/voting'

import MovieView from '../views/movie'

export default PageView.extend({
    pageTitle: 'Movie Night',
    template: votingTemplate,
    events: {
        'click [data-hook~=action-vote]': 'vote',
        'click [data-hook~=action-refresh]': 'refresh',
        'click [data-hook~=action-suggest-movie]': 'suggestMovie',
    },
    render: function () {
        this.renderWithTemplate()
        this.renderCollection(this.collection,
            MovieView,
            this.queryByHook('movie-list'),
            {
                reverse: true
            });

        if (!this.collection.length) {
            this.fetchCollection();
        }
    },
    refresh() {
        this.fetchCollection()
    },
    vote() {
        let votes = this.collection.where({ selected: true })
            .map(movie => movie.id)

        $.ajax({
            url: '/api/polls/latest/vote',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                votes
            })
        })
        .then(() => {
            toastr.success('Ai votat cu succes. Merci!')
        })
        .fail(err => {
            toastr.error('S-a produs o eroare, incearca din nou. Sorry!')
        })
        .always(() => {
            this.fetchCollection();
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
                imdbURL: url
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
            this.fetchCollection();
        })
    },
    fetchCollection: function () {
        this.collection.fetch();
        return false;
    },
});
