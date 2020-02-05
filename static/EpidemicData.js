window.chartColors = [
    "rgb(255, 99, 132)",
    "rgb(255, 159, 64)",
    "rgb(255, 205, 86)",
    "rgb(75, 192, 192)",
    "rgb(54, 162, 235)",
    "rgb(153, 102, 255)",
    "rgb(201, 203, 207)"
];

function HttpGet(url, func) {
    var xml_http = new XMLHttpRequest();
    xml_http.open("GET", url);
    xml_http.send();
    xml_http.onreadystatechange = function () {
        if (xml_http.readyState == 4 && xml_http.status == 200) {
            func(xml_http.responseText);
        }
    };
}

function HandleResponse(responseText) {
    var data = JSON.parse(responseText);
    console.log(data)
    Add(data)
    Rate(data)
}

function Add(data) {
    var epidemic_data = {};
    var confirm_rate = {};
    var suspect_rate = {};
    epidemic_data["title"] = "新增人数";
    epidemic_data["lables"] = [];
    epidemic_data["dataset"] = [];
    confirm_rate["label"] = "确诊";
    confirm_rate["data"] = [];
    suspect_rate["label"] = "疑似";
    suspect_rate["data"] = [];
    for (var i in data) {
        epidemic_data["lables"].push(data[i]["date"]);
        confirm_rate["data"].push(data[i]["confirm"]);
        suspect_rate["data"].push(data[i]["suspect"]);
    }
    epidemic_data["dataset"].push(confirm_rate);
    epidemic_data["dataset"].push(suspect_rate);
    var ctx = document.getElementById("add").getContext("2d");
    MakeChart(ctx, epidemic_data);
}

function Rate(data) {
    var epidemic_data = {};
    var confirm_rate = {};
    var suspect_rate = {};
    epidemic_data["title"] = "新增人数的速率";
    epidemic_data["lables"] = [];
    epidemic_data["dataset"] = [];
    confirm_rate["label"] = "确诊";
    confirm_rate["data"] = [];
    suspect_rate["label"] = "疑似";
    suspect_rate["data"] = [];
    for (var i in data) {
        epidemic_data["lables"].push(data[i]["date"]);
        confirm_rate["data"].push(data[i]["confirm_rate"]);
        suspect_rate["data"].push(data[i]["suspect_rate"]);
    }
    epidemic_data["dataset"].push(confirm_rate);
    epidemic_data["dataset"].push(suspect_rate);
    var ctx = document.getElementById("rate").getContext("2d");
    MakeChart(ctx, epidemic_data);
}

function MakeChart(ctx, epidemic_data) {
    var config = {
        type: "line",
        data: {
            labels: epidemic_data["lables"],
            datasets: []
        },
        options: {
            responsive: true,
            title: {
                display: true,
                text: epidemic_data["title"]
            },
            tooltips: {
                mode: "index",
                intersect: false,
            },
            hover: {
                mode: "nearest",
                intersect: true
            },
            scales: {
                xAxes: [{
                    display: true,
                    scaleLabel: {
                        display: true,
                        labelString: "日期"
                    }
                }],
                yAxes: [{
                    display: true,
                    scaleLabel: {
                        display: true,
                        labelString: "人数"
                    }
                }]
            }
        }
    };
    for (var i in epidemic_data["dataset"]) {
        var temp = {
            label: epidemic_data["dataset"][i]["label"],
            backgroundColor: window.chartColors[i],
            borderColor: window.chartColors[i],
            data: epidemic_data["dataset"][i]["data"],
            fill: false,
        };
        config.data.datasets.push(temp);
    }
    var myChart = new Chart(ctx, config);
}
