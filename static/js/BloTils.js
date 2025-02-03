// BloTils: https://www.blotils.com
// This file is released under the ISC license: https://opensource.org/licenses/ISC

{
    'use strict';

    if (window.blotils && window.blotils.vars)  // Compatibility with very old version; do not use.
        window.blotils = window.blotils.vars
    else
        window.blotils = window.blotils || {}

    const count_url = "api/v1/count_like"

    let script_data_url = document.querySelector('script[data-blotils_url]')
    let script_data_blotils_like_btn = document.querySelector('script[data-blotils_like_btn]')
    if (script_data_url) {
        window.blotils.base_url = script_data_url.dataset.blotils_url
    }
    else {
        window.blotils.base_url = 'https://blotils.com'
    }

    if (script_data_blotils_like_btn) {
        window.blotils.like_button = document.getElementById(script_data_blotils_like_btn.dataset.blotils_like_btn)
    }
    else {
        window.blotils.like_button = document.getElementById("blotils_like_btn")
    }
    let build_url = function (uri, params) {
        let url = new URL(uri, window.blotils.base_url)
        for (let key in params) {
            url.searchParams.append(key, params[key])
        }
        return url

    }

    window.blotils.count_like = function () {
        const options = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                page: window.location.pathname,
            })
        };
        fetch(build_url(count_url), options)
            .then(response => response.json())
            .then(data => console.log(data))
            .catch(error => console.error(error));
    }

    window.blotils.get_likes = function () {
        fetch(build_url(count_url, { page: window.location.pathname }))
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                console.log(data);
            })
            .catch(error => console.error(error));
    }
    window.blotils.like_button.addEventListener("click", window.blotils.count_like);
    window.blotils.get_likes()
}
