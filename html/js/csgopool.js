
function showHelp(helpId) {

    help = document.getElementById(helpId)
    
    if(help) {
      help.style.display = "inherit"
      help.style.visibility = "visible"
    }
  
}

function closeHelp(helpId) {

    help = document.getElementById(helpId)
    
    if(help) {
      help.style.display = "none"
      help.style.visibility = "hidden"
    }
}