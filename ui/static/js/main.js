var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

// document.getElementById("createBit").addEventListener("submit", e => {
// 	e.preventDefault();

// 	const form = e.target
// 	const formData = {
// 		title: form.title.value,
// 		content: form.content.value,
// 		expires: parseInt(form.expires.value) 
// 	}

// 	const jsonData = JSON.stringify(formData)

// 	fetch("http://localhost:4000/bits/create", {
// 		method: "POST",
// 		headers: {
// 			"Content-Type": "application/json"
// 		},
// 		body: jsonData
// 	}).then(res => {
// 		console.log(res)
// 		if(res.redirected) {
// 			window.location.href = res.url
// 		}
// 	})
// })