const promise = new Promise((resolve, reject) => {
  const xhr = new XMLHttpRequest();
  xhr.open("GET", "https://api.npoint.io/61fdbc52c8ceb9a9d6ba", "true");
  console.log(xhr);

  xhr.onload = () => {
    if (xhr.status === 200) {
      resolve(JSON.parse(xhr.response));
    } else {
      reject("Error loading data.");
    }
  };
  xhr.onerror = () => {
    reject("Network disable.");
  };
  xhr.send();
});

async function getAllTestimonials() {
  const response = await promise;
  console.log(response);

  let testimonialHTML = "";
  response.forEach(function (item) {
    testimonialHTML += `
    <div
    class="card m-4 pe-2 pt-2 ps-2 col-xs-2 col-sm-6 col-md-3 col-lg-2"
    style="width:300px; height:340px;"
  >
    <img
      src="${item.image}"
      alt="irsyad-photo"
      class="card-img-top p-1"
      style="width: 280px; height: 50%; object-fit: cover;"
    />
    <div class="card-body">
      <div class="card-text">
        <p class="qoute fw-bolder">${item.qoute}</p>
        <p class="author text-end fst-italic">- ${item.author}</p>
        <p class="text-end">${item.rating} <i class="fa-solid fa-star"></i></p>
      </div>
    </div>
  </div>
                        `;
  });

  document.getElementById("testimonials").innerHTML = testimonialHTML;
}

getAllTestimonials();

async function getFilteredTestimonials(rating) {
  const response = await promise;

  const testimonialsfiltered = response.filter((item) => {
    return item.rating === rating;
  });

  let testimonialHTML = "";

  if (testimonialsfiltered.length === 0) {
    testimonialHTML = "<h1>Data not found!</h1>";
  } else {
    testimonialsfiltered.forEach((item) => {
      testimonialHTML += `
      <div
      class="card m-4 pe-2 pt-2 ps-2 col-xs-2 col-sm-6 col-md-3 col-lg-2"
      style="width:300px; height:340px;"
    >
      <img
        src="${item.image}"
        alt="irsyad-photo"
        class="card-img-top p-1"
        style="width: 280px; height: 50%; object-fit: cover;"
      />
      <div class="card-body">
        <div class="card-text">
          <p class="qoute fw-bolder">${item.qoute}</p>
          <p class="author text-end fst-italic">- ${item.author}</p>
          <p class="text-end">${item.rating} <i class="fa-solid fa-star"></i></p>
        </div>
      </div>
    </div>
        `;
    });
  }

  document.getElementById("testimonials").innerHTML = testimonialHTML;
}
