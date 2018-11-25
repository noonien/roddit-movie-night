import View from 'ampersand-view'
import movieTemplate from '../../templates/includes/movie'

export default View.extend({
    template: movieTemplate,
    bindings: {
        'model.title': {
            hook: 'name'
        },
        'model.poster': {
            type: (el, val) => {
                console.log('el, val:', el, val);
                
                $(el).css('background-image', `url(${val})`)
            },
            hook: 'poster',
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
            hook: 'votes',
        },
        'model.plot': {
            hook: 'plot',
        },
        'model.genre': {
            hook: 'genre',
        },
        'model.year': {
            hook: 'year',
        },
        'model.rating_imdb': {
            hook: 'rating',
        },
    },
    events: {
        'click [data-hook~=action-select]': 'select',
    },
    select() {
        this.model.selected = !this.model.selected
        
        this.parent.vote(this.model.id, this.model.selected)
    },
});
