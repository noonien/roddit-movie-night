import AmpersandModel from 'ampersand-model'


export default AmpersandModel.extend({
    props: {
        id: 'any',
        votes: ['number', true, 0],
        title: ['string', true, ''],
        year: ['string', true, ''],
        genre: ['string', true, ''],
        plot: ['string', true, ''],
        poster: ['string', true, ''],
        imdb_url: ['string', true, ''],
        ratings: 'array'
    },
    session: {
        selected: ['boolean', true, false]
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
