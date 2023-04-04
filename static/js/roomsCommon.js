function RoomCheckAvailabilityPrompt(id, csrftoken) {
    let html = `
            <form id="check-availability-form" action="" method ="post" novalidate class="needs-validation">
                <div class="form-row">
                    <div class="col">
                        <div class="row" id="reservation-dates-modal">
                            <div class="col">
                                <input disabled required class="form-control" type="text" id="start" placeholder="Arrival" name="start">
                            </div>
                        
                            <div class="col">
                                <input disabled required class="form-control" type="text" id="end" placeholder="Departure" name="end">
                            </div> 
                        </div>
                    </div>
                </div>
            </form>
            `;
    
    attention.custom({
    msg: html, 
    title: "Choose your dates",
    willOpen: () => {
        const elem = document.getElementById("reservation-dates-modal");
        const rp = new DateRangePicker(elem, {
            format: "yyyy-mm-dd",
            showOnFocus: true,
            minDate: new Date(),
        });
    },

    didOpen: () => {
        document.getElementById("start").removeAttribute("disabled");
        document.getElementById("end").removeAttribute("disabled");
    },

    callback: function(result) {

        let form = document.getElementById("check-availability-form");
        let formData = new FormData(form);
        formData.append("csrf_token", csrftoken);
        formData.append("room_id", id);

        fetch('/search-availability-json', {
            method: "post",
            body: formData,
        })
            .then(response => response.json())
            .then(data => {
                if(data.ok) {
                attention.custom({
                    icon: 'success',
                    showConfirmButton: false,
                    msg: '<p>Room is available!</p>'
                        + '<p><a href="/book-room?id='
                        + data.room_id
                        + '&s='
                        + data.start_date
                        + '&e='
                        + data.end_date
                        + '" class="btn btn-primary">'
                        + 'Book now!</a></p>',
                })
                } else{
                attention.error({
                    msg: "No Availability",
                })
                }
                
            })
    }
    })
}