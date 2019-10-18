// ./app.js

new Vue({
    el: '#app',
    data: {
        sid: '',
        msg: '',
        resultQuest: [],
        resultChar: [],
    },
    methods: {
        login() {
            {
                var uid = document.getElementById("uid").value
                console.log("uid = " + uid);
                axios
                    .get('http://nt1.me:5000/login', { params: { uid: uid } })
                    .then(response => {
                        this.sid = response.data.data;
                        this.msg = "登入成功!";
                        console.log(this.sid)
                    })
                    .catch(function (error) {
                        console.log(error);
                    });
            }
        },
        status() {
            {
                this.msg = "查詢中...";
                var sid = document.getElementById("sid").textContent
                console.log("sid = " + sid);
                axios
                    .get('http://nt1.me:5000/status', { params: { sid: sid } })
                    .then(response => {
                        var data = JSON.stringify(response.data.data, null, 4)
                        console.log(data);
                        this.msg = data;
                    })
                    .catch(function (error) {
                        console.log(error);
                    });
            }
        },
        queryQuest(){
            console.log("enter queryQuest()")
            var questName = document.getElementById("queryText").value;
            console.log("Quest name to query = " + questName)
            axios.get("http://nt1.me:5000/query_quest", { params: { name: questName}})
                .then(response => { 
                    console.log(response.data.data)
                    this.resultQuest = response.data.data;
                })
        },
        playQuest(qtype, qid) {
            console.log("enter playQuest(), questID = " + qid);
            var sid = document.getElementById("sid").textContent
            axios.get("http://nt1.me:5000/play_quest", { params: {
                sid: sid,
                qtype: qtype,
                qid: qid,
                pt: 0
            } })
                .then(response => {
                    console.log(response.data.message)
                    window.alert(response.data.message);
                }).catch(function (error) {
                    window.alert(response.data.message);
                    console.log(error);
                });
        },
        queryChar() {
            console.log("enter queryChar()")
            var char = document.getElementById("queryText").value;
            console.log("Quest name to query = " + char)
            axios.get("http://nt1.me:5000/char", { params: { name: char } })
                .then(response => {
                    console.log(response.data.data)
                    this.resultChar = response.data.data;
                })
        },
        prompt() {
            var sid = document.getElementById("sid").textContent;
            this.msg = sid;
            console.log(sid);
        },
        prompt2() {
            this.msg = '快來看這裡！我是新訊息2！';
        }
    }
});
