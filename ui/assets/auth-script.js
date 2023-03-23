const sid = crypto.randomUUID()
window.onload = function (e) {
    selectFirstUser();
    document.getElementById('nav-claims-tab').addEventListener('click', generateIdTokenClaims);
}

function selectFirstUser() {
    document.getElementById('user-selection-0').click();
}

function fillAuthenticationForm(user) {
    document.getElementById('subject').value = user.sub;
    document.getElementById('given_name').value = user.given_name;
    document.getElementById('family_name').value = user.family_name;
    document.getElementById('birthdate').value = user.birthdate;

    const phone = document.getElementById('phone_number');
    if (phone) {
        phone.value = user.phone_number;
    }

    changeRadioSelection('acr', user.acr);
    changeRadioSelection('amr', user.amr);
    generateIdTokenClaims();
}

function generateIdTokenClaims() {
    const issuedAt = Math.floor(Date.now() / 1000);
    const formElem = document.getElementById('auth_form');
    const formData = new FormData(formElem);
    const claims = Object.fromEntries(formData.entries())
    claims.aud = formData.getAll('aud');
    claims.amr = formData.getAll('amr');
    claims.sid = sid;
    claims.jti = crypto.randomUUID();
    claims.auth_time = issuedAt;
    claims.iat = issuedAt;
    claims.exp = issuedAt + 900;
    claims.rat = issuedAt;
    claims.at_hash = 'YjllNTNjMDY1Y2MxZGNlYmQwODZiZDQwZDkzNzRjNGNjZDQ3YWFlMjgzN2IwZTQ1NTcxODlhMTU4NzhiOWE4Nw=='
    document.getElementById('id_token_claims').value = JSON.stringify(claims, Object.keys(claims).sort(), 2);
}

function generateLogoutTokenClaims(checkbox) {
    const claimsElem = document.getElementById('logout_token_claims')
    if (checkbox.checked) {
        claimsElem.removeAttribute('disabled');
        claimsElem.classList.remove('d-none');
        const issuedAt = Math.floor(Date.now() / 1000);
        const formElem = document.getElementById('auth_form');
        const formData = new FormData(formElem);
        const claims = {
            aud: formData.getAll('aud'),
            events: {
                'http://schemas.openid.net/event/backchannel-logout': {}
            },
            iat: issuedAt,
            iss: formData.get("iss"),
            jti: crypto.randomUUID(),
            sid: sid
        }
        claimsElem.value = JSON.stringify(claims, null, 2);
    } else {
        claimsElem.classList.add('d-none');
        claimsElem.setAttribute('disabled', 'disabled');
    }
}

function changeRadioSelection(radioSelectionName, valueToSelect) {
    let radioSelections = document.getElementsByName(radioSelectionName);
    for (let i = 0, length = radioSelections.length; i < length; i++) {
        if (radioSelections[i].value === valueToSelect) {
            radioSelections[i].checked = true;
            break;
        }
    }
}
