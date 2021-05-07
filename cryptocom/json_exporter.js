(function () {
  var date = new Date();
  var dateStr = date.getFullYear() + "-" + ("00" + (date.getMonth() + 1))
    .slice(-2) + "-" + ("00" + date.getDate()).slice(-2) + " " + ("00" + date
      .getHours()).slice(-2) + ":" + ("00" + date.getMinutes()).slice(-2) +
    ":" + ("00" + date.getSeconds()).slice(-2);

  function withdrawals() {
    return $.ajax({
        url: "https://crypto.com/fe-ex-api/record/withdraw_list",
        type: "POST",
        dataType: "json",
        contentType: "application/json;charset=utf-8",
        headers: {
          "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
        },
        data: JSON.stringify({
          "uaTime": dateStr,
          "securityInfo": "{\"timestamp\":\"" + dateStr + "\",\"meta\":{}}",
          "pageSize": 200,
          "page": 1,
          "coinSymbol": null
        })
      });
  }

  function deposits() {
    return $.ajax({
      url: "https://crypto.com/fe-ex-api/record/deposit_list",
      type: "POST",
      dataType: "json",
      contentType: "application/json;charset=utf-8",
      headers: {
        "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
      },
      data: JSON.stringify({
        "uaTime": dateStr,
        "securityInfo": "{\"timestamp\":\"" + dateStr + "\",\"meta\":{}}",
        "pageSize": 200,
        "page": 1,
        "coinSymbol": null
      })
    });
  }

  function crostaking() {
    return $.ajax({
      url: "https://crypto.com/fe-ex-api/record/staking_interest_history",
      type: "POST",
      dataType: "json",
      contentType: "application/json;charset=utf-8",
      headers: {
        "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
      },
      data: JSON.stringify({
        "uaTime": dateStr,
        "securityInfo": "{\"timestamp\":\"" + dateStr +
          "\",\"meta\":{}}",
        "pageSize": 200,
        "page": 1
      })
    });
  }

  function softstaking() {
    return $.ajax({
        url: "https://crypto.com/fe-ex-api/record/soft_staking_interest_list",
        type: "POST",
        dataType: "json",
        contentType: "application/json;charset=utf-8",
        headers: {
          "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
        },
        data: JSON.stringify({
          "uaTime": dateStr,
          "securityInfo": "{\"timestamp\":\"" + dateStr +
            "\",\"meta\":{}}",
          "pageSize": 200,
          "page": 1
        })
      });
  }

  function rebates() {
    return $.ajax({
      url: "https://crypto.com/fe-ex-api/record/rebate_trading_fee_history",
      type: "POST",
      dataType: "json",
      contentType: "application/json;charset=utf-8",
      headers: {
        "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
      },
      data: JSON.stringify({
        "uaTime": dateStr,
        "securityInfo": "{\"timestamp\":\"" + dateStr +
          "\",\"meta\":{}}",
        "pageSize": 200,
        "page": 1
      })
    });
  }

  function syndicates() {
    return $.ajax({
      url: "https://crypto.com/fe-ex-api/syndicate/user/activities?isCompleted=true&page=1&pageSize=10",
      type: "GET",
      dataType: "json",
      contentType: "application/json;charset=utf-8",
      headers: {
        "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
      }
    });
  }

  function supercharger() {
    return $.ajax({
      url: "https://crypto.com/fe-ex-api/record/supercharger_reward_history",
      type: "POST",
      dataType: "json",
      contentType: "application/json;charset=utf-8",
      headers: {
        "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
      },
      data: JSON.stringify({
        "uaTime": dateStr,
        "securityInfo": "{\"timestamp\":\"" + dateStr + "\",\"meta\":{}}",
        "pageSize": 200,
        "page": 1
      })
    });
  }

  function referralbonus() {
    return $.ajax({
      url: "https://crypto.com/fe-ex-api/referral/bonus/history?page=1&pageSize=200",
      type: "GET",
      dataType: "json",
      contentType: "application/json;charset=utf-8",
      headers: {
        "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
      }
    });
  }

  function referraltradecommission() {
    return $.ajax({
      url: "https://crypto.com/fe-ex-api/referral/trade_commission/history?page=1&pageSize=200",
      type: "GET",
      dataType: "json",
      contentType: "application/json;charset=utf-8",
      headers: {
        "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
      }
    });
  }

  function referralreward() {
    return $.ajax({
      url: "https://crypto.com/fe-ex-api/referral/reward/info",
      type: "GET",
      dataType: "json",
      contentType: "application/json;charset=utf-8",
      headers: {
        "exchange-token": document.cookie.match(/token=([0-9a-zA-Z]+)/)[1]
      }
    });
  }

  $.when(withdrawals(), deposits(), crostaking(), softstaking(), rebates(), syndicates(), supercharger(), referralbonus(), referraltradecommission(), referralreward()).done(
      function(withs, deps, cros, sstake, rebs, syn, sup, bon, tcom, rew){
    var j = "{";
    if (withs[2].status == 200) {
      j += "\"withs\":"+JSON.stringify(withs[0].data)+",";
    }
    if (deps[2].status == 200) {
      j += "\"deps\":"+JSON.stringify(deps[0].data)+",";
    }
    if (cros[2].status == 200) {
      j += "\"cros\":"+JSON.stringify(cros[0].data)+",";
    }
    if (sstake[2].status == 200) {
      j += "\"sstake\":"+JSON.stringify(sstake[0].data)+",";
    }
    if (rebs[2].status == 200) {
      j += "\"rebs\":"+JSON.stringify(rebs[0].data)+",";
    }
    if (syn[2].status == 200) {
      j += "\"syn\":"+JSON.stringify(syn[0].data)+",";
    }
    if (sup[2].status == 200) {
      j += "\"sup\":"+JSON.stringify(sup[0].data)+",";
    }
    if (tcom[2].status == 200) {
      j += "\"tcom\":"+JSON.stringify(tcom[0].data)+",";
    }
    if (bon[2].status == 200) {
      j += "\"bon\":"+JSON.stringify(bon[0].data)+",";
    }
    if (rew[2].status == 200) {
      j += "\"rew\":"+JSON.stringify(rew[0].data)+",";
    }
    if (j.length > 1) {
      j = j.slice(0, -1)
    }
    j += "}"
    // Download the JSON
    let o = encodeURI("data:application/json;charset=utf-8," + j), link = document.createElement("a");
    link.setAttribute("href", o), link.setAttribute("download", "CdC_Ex_ExportJS.json"),
      document.body.appendChild(link), link.click();
  });
})();
