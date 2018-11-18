import View from 'ampersand-view'
import movieTemplate from '../../templates/includes/movie'

export default View.extend({
    template: movieTemplate,
    bindings: {
        'model.name': '[data-hook~=name]',
        'model.image': {
            type: 'attribute',
            hook: 'image',
            name: 'src'
        },
        'model.imdbURL': {
            type: 'attribute',
            hook: 'name',
            name: 'href'
        },
        'model.selected': {
            type: 'booleanAttribute',
            hook: 'selected',
            name: 'checked'
        },
        'model.numVotes': {
            type: 'text',
            hook: 'votes',
        },
    },
    events: {
        'click [data-hook~=action-select]': 'select',
    },
    select() {
        this.model.selected = !this.model.selected
    },
});
