// ./app.js

new Vue({
    el: '#app',
    data: {
        sid: 'N/A',
        msg: '',
        resultQuest: [],
        resultChar: [],
        resultUzu: [],
        resultGacha: [],
        resultGachaInfo: [],
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
            var sid = document.getElementById("sid").textContent;
            if (sid==""){
                window.alert("請先登入");
                return;
            }
            console.log("enter playQuest(), questID = " + qid);            
            axios.get("http://nt1.me:5000/play_quest", { params: {
                sid: sid,
                qtype: qtype,
                qid: qid,
                pt: 0
            } })
                .then(response => {
                    console.log(response.data.status);
                    if (response.data.status != 200){
                        window.alert("failed");
                    }else{
                        window.alert(response.data.message);
                    }
                }).catch(function (error) {
                    window.alert(error);
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
        queryUzu() {
            console.log("enter queryUzu()")
            axios.get("http://nt1.me:5000/query_uzu")
                .then(response => {
                    console.log(response.data.data)
                    this.resultUzu = response.data.data;
                })
        },
        playUzu(scid, uzid) {
            var sid = document.getElementById("sid").textContent;
            if (sid == "") {
                window.alert("請先登入");
                return;
            }
            axios.get("http://nt1.me:5000/play_uzu", {
                params: {
                    sid: sid,
                    uzid: uzid,
                    scid: scid,
                }
            })
                .then(response => {
                    console.log(response.data.status);
                    if (response.data.status != 200) {
                        window.alert("failed");
                    } else {
                        window.alert(response.data.message);
                    }
                }).catch(function (error) {
                    window.alert(error);
                });
        },
        gacha() {
            var sid = document.getElementById("sid").textContent;
            var gachaID = document.getElementById("gachaID").value;
            var gachaCount = document.getElementById("gachaCount").value;
            if (sid == "") {
                window.alert("請先登入");
                return;
            }
            console.log(sid, gachaID, gachaCount);
            var tempResult = [];
            // for (index = 0; index < gachaBatch; index++) {                
            axios.get("http://nt1.me:5000/gacha", {
                params: {
                    sid: sid,
                    gacha_id: gachaID,
                    gacha_count: gachaCount,
                }
            })
                .then(response => {
                    console.log(response.data.status);
                    if (response.data.status != 200) {
                        window.alert("failed");
                        console.log(response.data);
                    } else {
                        console.log(response.data.data[0].Name);
                        // tempResult.push(response.data.data);
                        this.resultGacha = response.data.data
                    }
                }).catch(function (error) {
                    window.alert(error);
                });
                // console.log(tempResult);
                // this.resultGacha = tempResult;
            // }
        },
        queryGachaInfo() {
            console.log("enter queryGachaInfo()")
            var sid = document.getElementById("sid").textContent;
            if (sid == "") {
                window.alert("請先登入");
                return;
            }
            axios.get("http://localhost:5000/events", {
                params: {
                    sid: sid,
                }
            }).then(response => {
                console.log(response.data.data)
                this.resultGachaInfo = response.data.data;
            })
        },
        info() {
            axios.defaults.headers.common['origin'] = "mysite.com";
            axios.get("http://localhost:5000/events?sid=f3fe67d003bb05789502f21e8c9dac0f", {
            })
                .then(response => {
                    console.log(response.data.status);
                }).catch(function (error) {
                    window.alert(error);
                });
            // console.log(tempResult);
            // this.resultGacha = tempResult;
            // }
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
