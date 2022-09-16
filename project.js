let dataProject = []
function addProject(event){
    event.preventDefault()

    let title = document.getElementById("project-name").value
    let content = document.getElementById("desc").value
    let startDate = document.getElementById("start-date").value
    let endDate = document.getElementById("end-date").value
    let node = document.getElementById("cb1").checked
    let next = document.getElementById("cb2").checked
    let react = document.getElementById("cb3").checked
    let typescript = document.getElementById("cb4").checked
    let image = document.getElementById("input-upload-image").files

    // untuk membuat url gambar, agar tampil
    image = URL.createObjectURL(image[0])

    if(node){
        node = document.getElementById("cb1").value
    }else {
        node = ""
    }
    if(next){
        next = document.getElementById("cb2").value
    }else {
        next = ""
    }
    if(react){
        react = document.getElementById("cb3").value
    }else {
        react = ""
    }
    if(typescript){
        typescript = document.getElementById("cb4").value
    }else {
        typescript = ""
    }

    let Project = {
        title,
        content,
        startDate,
        endDate,
        node,
        next,
        react,
        typescript,
        duration: getDurationTime(new Date(startDate),new Date(endDate)),
        image,
    }

    dataProject.push(Project)
    console.log(dataProject);

   renderProject()
}

function renderProject(){
    
    document.getElementById("contents").innerHTML = ''

    console.log(dataProject);
    
    for (let index = 0; index < dataProject.length; index++) {
        
        // console.log(dataBlog[index]);
        document.getElementById("contents").innerHTML += 

        `<div id="contents" class="project-content">
        <div class="project-img">
            <img src=" ${dataProject[index].image}" alt="image">
            <a href="Project-detail.html">
            <h4>${dataProject[index].title}</h4></a>
            <p>Duration : ${dataProject[index].duration}</p>
        </div>
        <div class="content-fill">
            <p>${dataProject[index].content}</p>
        </div>
        <div class="i-tech">
            <i class="fa-brands fa-${dataProject[index].node}-js"></i>
            <i class="fa-brands fa-${dataProject[index].react}"></i>
            <i class="cib-${dataProject[index].next}-js"></i>
            <i class="cib-${dataProject[index].typescript}"></i> 
        </div>
        <div class="button-group">
            <button class="btn-edit">Edit</button>
            <button class="btn-delete">Delete</button>
        </div>

    </div>`
    }
}

function getDurationTime(start, end){

    let distance = end - start

    let milisecond = 1000 
    let secondInHours = 3600 
    let hoursInDay = 24 

    let distanceDay = Math.floor(distance / (milisecond * secondInHours * hoursInDay))

    let months = (end.getFullYear() - start.getFullYear()) * 12;
    months -= start.getMonth();
    months += end.getMonth();
    
    
    
    if(months > 0){
        return `${months} Month`
    } 
    else if(distanceDay > 0){
        return `${distanceDay} day`
    } 
}