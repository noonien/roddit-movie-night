import AmpersandModel from 'ampersand-model'


export default AmpersandModel.extend({
    props: {
        id: 'any',
        numVotes: ['number', true, 0],
        name: ['string', true, ''],
        image: ['string', true, ''],
        imdbURL: ['string', true, ''],
    },
    session: {
        selected: ['boolean', true, true]
    },
})
