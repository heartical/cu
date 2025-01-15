const go = new Go();
const roomID = "{{.}}";
fetch(`/static/vm.wasm`)
    .then(resp => resp.arrayBuffer())
    .then(bytes => bytes)
    .then(bytes => WebAssembly.instantiate(bytes, go.importObject))
    .then(function (result) {
        document.getElementById("loading-text").style.display = 'none';
        window.roomID = roomID;
        go.run(result.instance);
    })
    .catch(function (err) {
        console.error(err);
    });