var ul = document.querySelector('ul')

new EventSource('/events').addEventListener('message', function(e) {
  var message = e.data;

  ul.innerHTML += '<li>' + e.data + '</li>';
})
