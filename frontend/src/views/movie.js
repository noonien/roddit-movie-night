import View from 'ampersand-view'
import movieTemplate from '../../templates/includes/movie'

export default View.extend({
    template: movieTemplate,
    bindings: {
        'model.title': '[data-hook~=name]',
        'model.poster': {
            type: 'attribute',
            hook: 'image',
            name: 'src'
        },
        'model.imdb_url': {
            type: 'attribute',
            hook: 'name',
            name: 'href'
        },
        'model.selected': {
            type: 'booleanAttribute',
            hook: 'selected',
            name: 'checked'
        },
        'model.votes': {
            type: 'text',
            hook: 'votes',
        },
        'model.plot': {
            type: 'text',
            hook: 'plot',
        },
        'model.genre': {
            type: 'text',
            hook: 'genre',
        },
        'model.year': {
            type: 'text',
            hook: 'year',
        },
        'model.rating_imdb': {
            type: 'text',
            hook: 'rating',
        },
    },
    events: {
        'click [data-hook~=action-select]': 'select',
    },
    select() {
        this.model.selected = !this.model.selected
        
        this.parent.vote()
    },
});
