'use strict';

(function () {
    const match = document.getElementById('match');

    const matchSuccess = document.getElementById('matchSuccess');
    const matchFailure = document.getElementById('matchFailure');
    const matchError = document.getElementById('matchError');

    const matchTemplate = document.getElementById('matchTemplate');
    const matchURI = document.getElementById('matchURI');

    function executeMatch() {
        const currentURL = new URL(document.location.href);
        currentURL.searchParams.set('matchTemplate', matchTemplate.value);
        currentURL.searchParams.set('matchURI', matchURI.value);
        history.replaceState({}, '', currentURL)

        const apiURL = new URL('/match', document.location.href);
        apiURL.searchParams.append('template', matchTemplate.value)
        apiURL.searchParams.append('uri', matchURI.value)

        fetch(apiURL)
            .then(data => data.json())
            .then(data => {
                if (data.Match) {
                    matchSuccess.querySelector('pre code').textContent = JSON.stringify(data.Values, null, 2);
                    matchFailure.classList.remove('show');
                    matchError.classList.remove('show');
                    matchSuccess.classList.add('show');

                    return;
                }

                if (false === data.Match) {
                    matchSuccess.classList.remove('show');
                    matchError.classList.remove('show');
                    matchFailure.classList.add('show');

                    return;
                }

                if (data.Error) {
                    matchFailure.classList.remove('show');
                    matchSuccess.classList.remove('show');
                    matchError.querySelector('h4').textContent = data.Error;
                    matchError.classList.add('show');
                }
            });
    }

    match.onsubmit = function (e) {
        e.preventDefault();
        executeMatch();
    };

    const urlParams = new URLSearchParams(window.location.search);

    const matchTemplateParam = urlParams.get('matchTemplate');
    if (matchTemplateParam) matchTemplate.value = matchTemplateParam;

    const matchURIParam = urlParams.get('matchURI');
    if (matchURIParam) matchURI.value = matchURIParam;

    if (matchURIParam && matchURIParam) executeMatch();
})();
