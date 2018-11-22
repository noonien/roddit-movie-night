import AmpersandModel from 'ampersand-model'


export default AmpersandModel.extend({
    props: {
        id: 'any',
        votes: 'number',
        title: 'string',
        year: 'string',
        genre: 'string',
        plot: 'string',
        poster: 'string',
        imdb_url: 'string',
        ratings: 'array',
    },
    session: {
        selected: 'boolean',
    },
    derived: {
        rating_imdb: {
            deps: ['ratings'],
            fn: function () {
                return this.ratings.find(rat => rat.source === 'Internet Movie Database').rating
            }
        },
    },
    parse(resp) {
        resp.imdb_url = resp.url
        delete resp['url']

        return resp 
    }
})
