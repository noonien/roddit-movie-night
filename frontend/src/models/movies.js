import Collection from 'ampersand-rest-collection'
import Movie from './movie'

export default Collection.extend({
  model: Movie,
  url: '/api/polls/latest/movies',
  comparator: 'numVotes',
})
