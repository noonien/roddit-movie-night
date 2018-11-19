import AmpersandModel from 'ampersand-model'
import Movies from './movies'


export default AmpersandModel.extend({
    url: '/api/polls/latest',
    props: {
        poll: {
            id: 'any',
            name: ['string', true, ''],
            info: ['string', true, ''],
            created_at: ['string', true, ''],
            updated_at: ['string', true, ''],
            closes_at: ['string', true, ''],
        },
        votes: 'array'
    },
    collections: {
        movies: Movies
    }
})
