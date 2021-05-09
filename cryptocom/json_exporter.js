(function () {
  var date = new Date();
  var dateStr = date.toISOString().split(/[T\.]/).slice(0, -1).join(" ");
  function getData(method, endpoint, additionalKeys = {}) {
    if (method == "POST") {
      var data = {
        uaTime: dateStr,
        securityInfo: { timestamp: dateStr, meta: {} },
        pageSize: 200,
        page: 1,
      };
      if (additionalKeys)
        for (const [key, value] of Object.entries(additionalKeys)) {
          data[key] = value;
        }
      return $.ajax({
        url: "https://crypto.com/fe-ex-api/" + endpoint,
        type: "POST",
        dataType: "json",
        contentType: "application/json;charset=utf-8",
        headers: {
          "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1],
        },
        data: JSON.stringify(data),
      });
    } else if (method == "GET") {
      return $.ajax({
        url: "https://crypto.com/fe-ex-api/" + endpoint,
        type: "GET",
        dataType: "json",
        contentType: "application/json;charset=utf-8",
        headers: {
          "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1],
        },
      });
    }
  }

  $.when(
    getData("POST", "record/withdraw_list", { coinSymbol: null }),
    getData("POST", "record/deposit_list", { coinSymbol: null }),
    getData("POST", "record/staking_interest_history"),
    getData("POST", "record/soft_staking_interest_list"),
    getData("POST", "record/rebate_trading_fee_history"),
    getData(
      "GET",
      "syndicate/user/activities?isCompleted=true&page=1&pageSize=10"
    ),
    getData("POST", "record/supercharger_reward_history"),
    getData("GET", "referral/bonus/history?page=1&pageSize=200"),
    getData("GET", "referral/trade_commission/history?page=1&pageSize=200"),
    getData("GET", "referral/reward/info")
  ).done(function (withs, deps, cros, sstake, rebs, syn, sup, bon, tcom, rew) {
    var j = {};
    if (withs[2].status == 200) j["withs"] = withs[0].data;
    if (deps[2].status == 200) j["deps"] = deps[0].data;
    if (cros[2].status == 200) j["cros"] = cros[0].data;
    if (sstake[2].status == 200) j["sstake"] = sstake[0].data;
    if (rebs[2].status == 200) j["rebs"] = rebs[0].data;
    if (syn[2].status == 200) j["syn"] = syn[0].data;
    if (sup[2].status == 200) j["sup"] = sup[0].data;
    if (tcom[2].status == 200) j["tcom"] = tcom[0].data;
    if (bon[2].status == 200) j["bon"] = bon[0].data;
    if (rew[2].status == 200) j["rew"] = rew[0].data;
    // Download the JSON
    var o = encodeURI("data:text/json;charset=utf-8," + JSON.stringify(j));
    var link = document.createElement("a");
    link.setAttribute("href", o);
    link.setAttribute("download", "CdC_Ex_ExportJS.json");
    document.body.appendChild(link);
    link.click();
    link.remove();
  });
})();
