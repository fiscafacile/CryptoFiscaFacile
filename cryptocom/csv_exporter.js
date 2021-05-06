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
        contentType: "application/json",
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
      contentType: "application/json",
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
      contentType: "application/json",
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
        contentType: "application/json",
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
      contentType: "application/json",
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
      contentType: "application/json",
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
      contentType: "application/json",
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

  $.when(withdrawals(), deposits(), crostaking(), softstaking(), rebates(), syndicates(), supercharger()).done(
      function(withs, deps, cros, stake, rebs, syn, sup){
    var t = ["Date", "Sent Amount", "Sent Currency", "Received Amount",
      "Received Currency", "Fee Amount", "Fee Currency", "Net Worth Amount",
      "Net Worth Currency", "Label", "Description", "TxHash"].join(",");

    if (withs[2].status == 200) {
      withs[0].data.financeList.forEach(function(e) {
        t += "\n" + [new Date(parseInt(e.updateAtTime)).toISOString(),
          e.amount, e.symbol, "", "", e.fee, e.symbol, "", "", "",
          "Withdrawal to " + e.addressTo + " (" + e.status_text +
          ")", e.txid].join(",");
      });
    }

    if (deps[2].status == 200) {
      deps[0].data.financeList.forEach(function(e) {
        t += "\n" + [new Date(parseInt(e.updateAtTime)).toISOString(),
          "", "", e.amount, e.symbol, "", "", "", "", "",
          "Deposit from " + e.addressTo + " (" + e.status_text +
          ")", e.txid].join(",");
      });
    }

    if (cros[2].status == 200) {
      cros[0].data.historyList.forEach(function(e) {
        t += "\n" + [new Date(parseInt(e.createdAtTime)).toISOString(),
          "", "", e.interestAmount, e.coinSymbol, "", "", "", "", "Reward",
          "Interest on " + e.stakeAmount + " at " + e.apr * 100 + "% APR (" + e.status_text + ")",
          ""].join(",");
      });
    }

    if (stake[2].status == 200) {
      stake[0].data.softStakingInterestList.forEach(function(e) {
        t += "\n" + [new Date(parseInt(e.mtime)).toISOString(), "",
          "", e.amount, e.coinSymbol, "", "", "", "", "Reward",
          "Interest on " + e.principal + " " + e.coinSymbol +
          " at " + e.apr * 100 + "% APR", ""].join(",");
      });
    }

    if (rebs[2].status == 200) {
      rebs[0].data.historyList.forEach(function(e) {
        t += "\n" + [new Date(parseInt(e.createdAtTime))
        .toISOString(), "", "", e.rebateAmount, e.coinSymbol, "",
          "", "", "", "Reward", "Rebate on " + e.feePaid + " " + e
          .coinSymbol + " at " + e.rebatePercentage * 100 + "%", ""].join(",");
      });
    }

    if (syn[2].status == 200) {
      syn[0].data.activities.forEach(function(e) {
        t += "\n" + [new Date(parseInt(e.userModifyTime))
          .toISOString(), e.committedCRO - e.refundedCRO, "CRO", e
          .allocatedVolume, e.syndicateCoin, "", "", "", "",
          "Syndicate", e.syndicateCoin + " syndicate at " + e
          .discountRate * 100 + "% off (" + e.discountedPrice +
          "CRO)", ""].join(",");
      });
    }

    if (sup[2].status == 200) {
      sup[0].data.historyList.forEach(function(e) {
        t += "\n" + [new Date(parseInt(e.createdAt)).toISOString(),
          "", "", e.rewardAmount, e.coinSymbol, "", "", "", "",
          "Reward", e.coinSymbol + " Supercharger reward", ""].join(",");
      });
    }

    // Remove transfers to/from the app
    t = t.replace(/^.*Crypto.com\sApp.*($|\r\n|\r|\n)/gm, "");

    // Remove blank lines from the output
    t = t.replace(/^\s*[\r\n]/gm, "")

    // Download the CSV
    let o = encodeURI("data:text/csv;charset=utf-8," + t), link = document.createElement("a");
    link.setAttribute("href", o), link.setAttribute("download", "crypto_exchange_data.csv"),
      document.body.appendChild(link), link.click();
  });
})();
