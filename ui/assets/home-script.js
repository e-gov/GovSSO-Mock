window.onload = function(e){
    selectFirstClient();
}

function selectFirstClient() {
    document.getElementById("client-selection-0").click();
}

function fillBackchannelLogoutForm(client) {
    document.getElementById("backchannel_logout_uri").value = client.backchannel_logout_uri;
    document.getElementById("client_id").value = client.client_id;
}