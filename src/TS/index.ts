window.onscroll = function () {
  scrollResponse();
};

function scrollResponse() {
  var navbar = document.getElementById("navbar");
  if (navbar === null) {
    throw new Error("navbar is null");
  }
  var sticky = navbar.offsetTop;

  if (sticky === undefined) {
    throw new Error("sticky is undefined");
  }
  if (window.scrollY >= sticky) {
    navbar?.classList.add("sticky");
  } else {
    navbar?.classList.remove("sticky");
  }
}
