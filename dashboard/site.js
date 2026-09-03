(function () {
  var stage = document.getElementById("stage");
  var burger = document.getElementById("burger");
  var menu = document.getElementById("menu");
  if (!stage || !burger || !menu) {
    return;
  }

  function setOpen(open) {
    stage.classList.toggle("is-open", open);
    burger.setAttribute("aria-expanded", String(open));
    burger.setAttribute("aria-label", open ? "Close menu" : "Open menu");
    menu.setAttribute("aria-hidden", String(!open));
    document.body.style.overflow = open ? "hidden" : "";
  }

  burger.addEventListener("click", function () {
    setOpen(!stage.classList.contains("is-open"));
  });

  menu.addEventListener("click", function (e) {
    if (e.target.closest("a")) {
      setOpen(false);
    }
  });

  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") {
      setOpen(false);
    }
  });

  window.addEventListener("resize", function () {
    if (window.innerWidth / window.innerHeight > 1.1) {
      setOpen(false);
    }
  });
})();
