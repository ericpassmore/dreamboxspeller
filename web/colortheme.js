// Toggle switch calls into this code to change colortheme
// sets the color theme accordingly by changing CSS classes
// note side effect
// preserve colortheme value in storage to presist across pages

// function to make CSS Color Theme Changes
// lastly set checkbox state to match theme
// acceptable values are "light" or "dark"
// all other values are simply ignored
function setcolor(theme) {
  if (theme.localeCompare("light") == 0) {
    $("body").removeClass("dark");
    $(".highlight-box").removeClass("dark");
    $("body").addClass("light");
    $(".highlight-box").addClass("light");
    $('#colortheme input[type="checkbox"]').prop("checked",false);
  }
  if (theme.localeCompare("dark") == 0) {
    $("body").removeClass("light");
    $(".highlight-box").removeClass("light");
    $("body").addClass("dark");
    $(".highlight-box").addClass("dark") ;
    $('#colortheme input[type="checkbox"]').prop("checked",true);
  }
}
// clear out storage values
// restore default theme
// should not need this but nice developer/test feature
function resetcolor() {
  localStorage.removeItem('colortheme')
  setcolor('dark')
}
$(document).ready(function(){
  // **** MAIN run after document is ready *****
  // first preserve value if it does not exist
  // then run function to set color
  if (!localStorage.getItem('colortheme')) {
    localStorage.setItem('colortheme','dark');
  } else {
    setcolor(localStorage.getItem('colortheme'));
  }
  // this only runs onclick for #colortheme checkbox
  // onclick sets colors and preserves value
  $('#colortheme input[type="checkbox"]').click(function(){
    if($(this).prop("checked") == false){
      console.log("Checkbox is checked.");
      setcolor("light");
      localStorage.setItem('colortheme','light');
    }
    else if($(this).prop("checked") == true){
      console.log("Checkbox is unchecked.");
      setcolor("dark");
      localStorage.setItem('colortheme','dark');
    }
  });
});
