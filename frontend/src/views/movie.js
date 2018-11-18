import View from 'ampersand-view'
import movieTemplate from '../../templates/includes/movie'

export default View.extend({
    template: movieTemplate,
    bindings: {
        'model.fullName': '[data-hook~=name]',
        'model.avatar': {
            type: 'attribute',
            hook: 'avatar',
            name: 'src'
        },
        'model.editUrl': {
            type: 'attribute',
            hook: 'action-edit',
            name: 'href'
        },
        'model.viewUrl': {
            type: 'attribute',
            hook: 'name',
            name: 'href'
        }
    },
    events: {
        'click [data-hook~=action-delete]': 'handleRemoveClick',
        'click [data-hook~=action-select]': 'select',
        // 'click [data-hook~=check]': 'ignoreClick',
    },
    ignoreClick(evt) {
        evt.preventDefault()
        return false
    },
    select() {
        let $check = this.queryByHook('checkbox');
        $check.checked = !$check.checked
        console.log('$check.checked:', $check.checked);
        
    },
    handleRemoveClick() {
        this.model.destroy();
        return false;
    }
});
