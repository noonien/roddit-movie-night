import Collection from 'ampersand-rest-collection'
import Movie from './movie'

export default Collection.extend({
  model: Movie,
  url: 'http://localhost:3000/api/polls/latest/movies',
  comparator: 'numVotes',
})
