window.onload = function (e) {
    selectFirstUser();
}

function selectFirstUser() {
    document.getElementById("user-selection-0").click();
}

function fillAuthenticationForm(user) {
    document.getElementById("subject").value = user.sub;
    document.getElementById("given_name").value = user.given_name;
    document.getElementById("family_name").value = user.family_name;
    document.getElementById("birthdate").value = user.birthdate;

    let phone = document.getElementById("phone")
    if (phone) {
        phone.value = user.phone;
    }

    changeRadioSelection("acr", user.acr);
    changeRadioSelection("amr", user.amr);
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
