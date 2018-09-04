var chatWindow = new Bubbles(document.getElementById("chat"), "chatWindow", {
    inputCallbackFn: function (o) {
        fetch("/api/",
            {
                headers: {
                    'Accept': 'application/json',
                    'matricula': '1004',
                    'Content-Type': 'application/json'
                },
                method: "POST",
                body: JSON.stringify(o.input)
            })
            .then(res => res.json())
            .then(function (res) {
                obj = {
                    ice: {
                        says: [res]
                    }
                }
                chatWindow.talk(obj)
            })
            .catch(function (error) { console.log(error) })
    }
});
var convo = {
    ice: {
        says: ["Olá", "Você gostria de marcar férias?"]
    }
}

chatWindow.talk(convo)
