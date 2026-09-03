(function () {
  "use strict";

  var TABS = ["opt-in", "unratified"];
  var PROVENANCE = {
    ratified: true,
    provisional: true,
    unofficial: true,
  };

  function $(sel, root) {
    return (root || document).querySelector(sel);
  }

  function escapeHtml(value) {
    return String(value == null ? "" : value)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;");
  }

  function formatGsi(gsi) {
    if (typeof gsi !== "number" || Number.isNaN(gsi)) {
      return "—";
    }
    return gsi.toFixed(1);
  }

  function formatGeneratedAt(iso) {
    if (!iso) {
      return "—";
    }
    var d = new Date(iso);
    if (Number.isNaN(d.getTime())) {
      return escapeHtml(iso);
    }
    return escapeHtml(d.toISOString().replace("T", " ").replace(/\.\d{3}Z$/, " Z"));
  }

  function provenanceBadge(provenance) {
    var key = PROVENANCE[provenance] ? provenance : "unofficial";
    var label = PROVENANCE[provenance] ? provenance : provenance || "unknown";
    return (
      '<span class="badge badge-' +
      escapeHtml(key) +
      '">' +
      escapeHtml(label) +
      "</span>"
    );
  }

  function sampleBadge(isSample) {
    if (!isSample) {
      return "";
    }
    return '<span class="badge badge-sample">sample</span>';
  }

  function flagBadges(entry) {
    var parts = [];
    var catastrophic = Array.isArray(entry.catastrophic) ? entry.catastrophic : [];
    var na = Array.isArray(entry.na) ? entry.na : [];

    if (catastrophic.length) {
      parts.push(
        '<span class="badge badge-catastrophic" title="' +
          escapeHtml(catastrophic.join(", ")) +
          '">catastrophic ×' +
          catastrophic.length +
          "</span>"
      );
    }
    if (na.length) {
      parts.push(
        '<span class="badge badge-unofficial" title="' +
          escapeHtml(na.join(", ")) +
          '">N/A ×' +
          na.length +
          "</span>"
      );
    }
    if (!parts.length) {
      return '<span class="muted">—</span>';
    }
    return '<div class="flags">' + parts.join("") + "</div>";
  }

  function adapterLabel(entry) {
    var name = entry.adapter || "—";
    var ver = entry.adapter_version;
    if (ver) {
      return escapeHtml(name) + " <span class=\"muted\">" + escapeHtml(ver) + "</span>";
    }
    return escapeHtml(name);
  }

  function rowHtml(entry) {
    return (
      "<tr>" +
      '<td><div class="framework-cell">' +
      "<span>" +
      escapeHtml(entry.framework || "—") +
      "</span>" +
      sampleBadge(entry.sample === true) +
      "</div></td>" +
      "<td>" +
      adapterLabel(entry) +
      "</td>" +
      "<td>" +
      escapeHtml(formatGsi(entry.gsi)) +
      "</td>" +
      "<td>" +
      escapeHtml(entry.grade || "—") +
      "</td>" +
      "<td>" +
      provenanceBadge(entry.provenance) +
      "</td>" +
      "<td>" +
      flagBadges(entry) +
      "</td>" +
      "<td>" +
      formatGeneratedAt(entry.generated_at) +
      "</td>" +
      "</tr>"
    );
  }

  function setStatus(tab, message, isError) {
    var wrap = document.querySelector('.table-wrap[data-tab="' + tab + '"]');
    if (!wrap) {
      return;
    }
    var status = $("[data-status]", wrap);
    if (!status) {
      return;
    }
    status.hidden = !message;
    status.textContent = message || "";
    status.classList.toggle("error", !!isError);
  }

  function renderTab(tab, entries) {
    var wrap = document.querySelector('.table-wrap[data-tab="' + tab + '"]');
    if (!wrap) {
      return;
    }
    var table = $("table", wrap);
    var tbody = $("tbody", wrap);
    var filtered = entries.filter(function (e) {
      return e && e.tab === tab;
    });

    if (!filtered.length) {
      tbody.innerHTML = "";
      table.hidden = true;
      setStatus(tab, "No entries yet.", false);
      return;
    }

    filtered.sort(function (a, b) {
      var ga = typeof a.gsi === "number" ? a.gsi : -1;
      var gb = typeof b.gsi === "number" ? b.gsi : -1;
      if (gb !== ga) {
        return gb - ga;
      }
      return String(a.framework || "").localeCompare(String(b.framework || ""));
    });

    tbody.innerHTML = filtered.map(rowHtml).join("");
    table.hidden = false;
    setStatus(tab, "", false);
  }

  function entryUrl(name) {
    // index.json lists bare filenames under ../data/
    if (/^https?:\/\//i.test(name) || name.indexOf("/") === 0) {
      return name;
    }
    return "../data/" + name.replace(/^\.\//, "");
  }

  function loadEntry(name) {
    return fetch(entryUrl(name)).then(function (res) {
      if (!res.ok) {
        throw new Error("Failed to load " + name + " (" + res.status + ")");
      }
      return res.json();
    });
  }

  function failAll(message) {
    TABS.forEach(function (tab) {
      setStatus(tab, message, true);
      var wrap = document.querySelector('.table-wrap[data-tab="' + tab + '"]');
      if (wrap) {
        var table = $("table", wrap);
        if (table) {
          table.hidden = true;
        }
      }
    });
  }

  fetch("../data/index.json")
    .then(function (res) {
      if (!res.ok) {
        throw new Error("Failed to load data/index.json (" + res.status + ")");
      }
      return res.json();
    })
    .then(function (index) {
      if (!Array.isArray(index)) {
        throw new Error("data/index.json must be a JSON array of entry filenames");
      }
      if (!index.length) {
        TABS.forEach(function (tab) {
          renderTab(tab, []);
        });
        return;
      }
      return Promise.all(index.map(loadEntry)).then(function (entries) {
        TABS.forEach(function (tab) {
          renderTab(tab, entries);
        });
      });
    })
    .catch(function (err) {
      failAll(err && err.message ? err.message : String(err));
    });
})();
