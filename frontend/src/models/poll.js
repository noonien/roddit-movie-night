import AmpersandModel from 'ampersand-model'
import Movies from './movies'


export default AmpersandModel.extend({
    url: '/api/polls/latest',
    props: {
        poll: {
            id: 'any',
            name: 'string',
            info: 'string',
            created_at: 'string',
            updated_at: 'string',
            closes_at: 'string',
        },
        votes: 'array'
    },
    collections: {
        movies: Movies
    }
})
