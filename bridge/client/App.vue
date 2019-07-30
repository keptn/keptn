<template>
  <div id="app">
    <b-navbar toggleable="lg" type="dark" sticky>
      <b-navbar-brand href="/">
        <span>
          <img class="logo-navbar" src="/logo.png">
        </span>
        keptn's bridge
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    </b-navbar>
    <div class="wrapper">
      <nav id="sidebar">
        <RootSearch/>
        <RootList/>
      </nav>
      <div class="container-fluid">
        <router-view/>
      </div>
    </div>
  </div>
</template>

<script>
import RootList from '@/components/RootList/index.vue';
import RootSearch from '@/components/RootSearch/index.vue';

export default {
  name: 'App',
  components: { RootList, RootSearch },
};
</script>


<style lang="less">
.navbar {
  background-color: #006bb8;
}

.logo-navbar {
  width: 2rem;
}

.wrapper {
  display: flex;
  align-items: stretch;
}

#sidebar {
  min-width: 450px;
  max-width: 450px;
  background: #eeeeee;
  transition: all 0.3s;
  min-height: calc(100vh - 50px);
}

#sidebar.active {
  margin-left: -250px;
}

#sidebar .sidebar-header {
  padding: 20px;
  background: #006bb8;
}

#sidebar ul.components {
  padding: 20px 0;
  border-bottom: 1px solid #47748b;
}
</style>
