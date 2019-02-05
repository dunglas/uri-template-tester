'use strict';

$(
    $('#match').submit(function (e) {
        e.preventDefault();
        const elError = $('#matchError');
        const elMatch = $('#matchSuccess');
        const elNoMatch = $('#matchFailure');
        elError.collapse('hide');
        elMatch.collapse('hide');
        elNoMatch.collapse('hide');

        const { elements: { template, uri } } = this;

        const u = new URL('/match', document.location.href);
        u.searchParams.append('template', template.value)
        u.searchParams.append('uri', uri.value)

        fetch(u)
            .then(data => data.json())
            .then(data => {
                if (data.Match) {
                    elMatch.find('pre code').text(JSON.stringify(data.Values, null, 2));
                    elMatch.collapse('show');

                    return;
                }

                if (false === data.Match) {
                    elNoMatch.collapse('show');

                    return;
                }

                if (data.Error) {
                    elError.find('h4').text(data.Error, null, 2);
                    elError.collapse('show');
                }
            });
    })
);
