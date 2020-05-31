(function () {
    const ppField = document.getElementById('password-protected');

    if (ppField) {
        ppField.checked = false;
    }

    const forms = document.getElementsByTagName('form');

    for (let i = 0; i < forms.length; i++) {
        forms[i].addEventListener("submit", onFormSubmit, false);
    }
})();

function togglePasswordContainer() {
    const el = document.getElementById('password-container');
    el.querySelector('#password').required = !el.classList.toggle('hidden');
}

function onFormSubmit() {
    const buttons = this.getElementsByClassName('cta-btn');

    for (let i = 0; i < buttons.length; i++) {
        buttons[i].getElementsByClassName('btn-text')[0].classList.add('hidden');
        buttons[i].getElementsByClassName('spinner')[0].classList.remove('hidden');
        buttons[i].disabled = true;
    }
}
