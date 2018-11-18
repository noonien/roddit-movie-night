import PageView from './base'
import votingTemplate from '../../templates/pages/voting'
import $ from 'jquery'

import MovieView from '../views/movie'

import { getOGData } from '../utils'


export default PageView.extend({
    pageTitle: 'Movie Night',
    template: votingTemplate,
    events: {
        'click [data-hook~=add]': 'addRandom',
        'click [data-hook~=action-vote]': 'vote',
        'click [data-hook~=action-suggest-movie]': 'suggestMovie',
    },
    render: function () {
        this.renderWithTemplate();
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
        })
        .always(() => {
            this.fetchCollection();
        })
    },
    fetchCollection: function () {
        this.collection.fetch();
        return false;
    },
    resetCollection: function () {
        this.collection.reset();
    },
    addRandom: function () {
        function getRandom(min, max) {
            return min + Math.floor(Math.random() * (max - min + 1));
        }
        var firstNames = 'Joe Harry Larry Sue Bob Rose Angela Tom Merle Joseph Josephine'.split(' ');
        var lastNames = 'Smith Jewel Barker Stephenson Rossum Crockford'.split(' ');

        console.log('this.collection:', this.collection);
        
        this.collection.create({
            id: getRandom(0, 6),
            numVotes: this.collection.length + 1,
            name: firstNames[getRandom(0, firstNames.length - 1)],
            image: 'https://m.media-amazon.com/images/M/MV5BMTkzMzgzMTc1OF5BMl5BanBnXkFtZTgwNjQ4MzE0NjM@._V1_UX182_CR0,0,182,268_AL_.jpg',
            selected: true,
            imdbURL: 'inb4'
        });
    }
});
