<template>
  <div>

    <div class="container">

      <div class="tab-pane fade show active">
          <div class="row">
              <div class="col-md-12">
                  <h3>Here's a list of your Invoices</h3>
                  <table class="table">
                  <thead>
                    <tr>
                      <th scope="col">Invoice #</th>
                      <th scope="col">Invoice Name</th>
                      <th scope="col">Status</th>
                      <th scope="col"></th>
                    </tr>
                  </thead>
                  <tbody>
                    <template v-for="invoice in invoices">
                      <tr>
                        <th scope="row">{{ invoice.id }}</th>
                        <td>{{ invoice.name }}</td>
                        <td v-if="invoice.paid == 0 "> Unpaid </td>
                        <td v-else> Paid </td>
                        <td ><router-link :to="{ name: 'SingleInvoice', params: { invoice_id: invoice.id }}" class="btn btn-success">TO INVOICE</router-link>  </td>
                      </tr>
                    </template>
                  </tbody>
                </table>
              </div>
          </div>
      </div>
    </div>
  </div>

</template>

<script>
import axios from "axios";

export default {
  name: "ViewInvoices",
  components: {},
  data() {
    return {
      invoices: [],
      user: '',
    };
  },
  mounted() {
    this.user = JSON.parse(localStorage.getItem('user'));
    axios
      .get(`http://localhost:3128/invoices/user/${this.user.id}`,
        {
          headers: {"Token": localStorage.getItem("token")}
        }
      )
      .then(res => {
        if (res.data.status == true) {
          this.invoices = res.data.invoices;
        }
      });
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1,
h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #ffffff;
}

.tab-pane {
  margin-top: 20px;
}
</style>
