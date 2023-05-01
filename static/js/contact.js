function contactSubmissionPrompt(csrftoken){
    attention.custom({
        callback: function(result) {
            let form = document.getElementById("contactForm");
            let formData = new FormData(form);
            formData.append("csrf_token", csrftoken);

            fetch('/submit-contact', {
                method: "post",
                body: formData,
            })
                .then(response => response.json())
                .then(data => {
                    if(data.ok) {
                        attention.custom({
                            icon: 'success',
                            showConfirmButton: false,
                            msg: "We have received your contact request!",
                        })
                    } else{
                        attention.error({
                            msg: "Contact Request was not made. Please contact admin",
                        })
                    }
                })
        }
    })
}