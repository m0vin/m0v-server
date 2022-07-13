{{define "menujs"}}
    <script type="text/javascript">
        function collapsenav() {
            $nav = document.getElementById("main-nav");
            $nav.classList.toggle("minimize-main-nav"); 
            $nav.classList.toggle("main-nav"); 
            $icon = document.getElementById("icon-arrow-left");
            $icon.classList.toggle("icon-arrow-left");
            $icon.classList.toggle("icon-arrow-right");
            $header = document.getElementById("header");
            $header.classList.toggle("header-minimize-main-nav"); 
        }
        function records() {

        }
        function inbox() {

        }
        function reports() {
            document.querySelector(".header-nav-link.icon-personal-records").classList.remove("active");
            document.querySelector(".header-nav-link.icon-inbox").classList.remove("active");
            document.querySelector(".header-nav-link.icon-gear").classList.remove("active");
            document.querySelector(".header-nav-link.icon-reports-bar").classList.toggle("active");
        }
        function gear() {
            document.querySelector(".header-nav-link.icon-personal-records").classList.remove("active");
            document.querySelector(".header-nav-link.icon-inbox").classList.remove("active");
            document.querySelector(".header-nav-link.icon-reports-bar").classList.remove("active");
            document.querySelector(".header-nav-link.icon-gear").classList.toggle("active");
        }
        document.getElementById("icon-arrow-left").addEventListener('click', collapsenav);
        document.querySelector(".icon-reports-bar").addEventListener('click', reports);
        document.querySelector(".icon-gear").addEventListener('click', gear);
        document.querySelector(".icon-personal-records").addEventListener('click', records);
        document.querySelector(".icon-inbox").addEventListener('click', inbox);
    </script>
{{end}}

