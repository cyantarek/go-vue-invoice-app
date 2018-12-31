<template>
  <div class="single-page">
    <Header v-bind:user="user"/>
    <!--  display invoice data -->
    <div class="invoice">
      <!-- display invoice name here -->
      <div class="container">
        <div class="row">
          <div class="col-md-12">
            <h3>Invoice #{{ invoice.id }} by {{ user.company_name }}</h3>
            <table class="table">
              <thead>
                <tr>
                  <th scope="col">#</th>
                  <th scope="col">Transaction Name</th>
                  <th scope="col">Price ($)</th>
                </tr>
              </thead>
              <tbody>
                <template v-for="txn in transactions">
                  <tr :key="txn.id">
                    <th>{{ txn.id }}</th>
                    <td>{{ txn.name }}</td>
                    <td>{{ txn.price }} </td>
                  </tr>
                </template>
              </tbody>
              <tfoot>
                <td></td>
                <td style="text-align: right">Total :</td>
                <td><strong>$ {{ total_price }}</strong></td>
              </tfoot>
            </table>
          </div>
        </div>

        <div class="row">
          <form @submit.prevent="send" class="col-md-12">
            <h3>Enter Recipient's Name and Email to Send Invoice</h3>
            <div class="form-group">
              <label for="">Recipient Name</label>
              <input type="text" required class="form-control" placeholder="eg Chris" v-model="recipient.name">
            </div>

            <div class="form-group">
              <label for="">Recipient Email</label>
              <input type="email" required placeholder="eg chris@invoiceapp.com" class="form-control" v-model="recipient.email">
            </div>

            <div class="form-group">
                <button class="btn btn-primary" >Send Invoice</button>
                {{ loading }}
                {{ status }}
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Header from "./Header";
import axios from "axios";

export default {
  name: "SingleInvoice",
  components: {
    Header
  },
  data() {
    return {
      invoice: {},
      transactions: [],
      user: JSON.parse(localStorage.getItem("user")),
      total_price: 0,
      recipient : {
        name: '',
        email: ''
      },
      loading : '',
      status: '',
    };
  },
  methods: {
    send() {
      // prepare the formdata
      this.loading = "Sending Invoice, please wait...."
      const formData = new FormData();
      formData.append("user", JSON.stringify(this.user));
      formData.append("recipient", JSON.stringify(this.recipient));
      axios.post("http://localhost:3128/sendmail", formData, {
        headers: {"x-access-token": localStorage.getItem("token")}
      }).then(res => {
        console.log(res);
        this.loading = '';
        this.status = 'Invoice Sent'
      });   
    }
  },
  mounted() {
    
    // make request to fetch invoice data
    this.user = JSON.parse(localStorage.getItem("user"));
    let token = localStorage.getItem("token");
    let invoice_id = this.$route.params.invoice_id;
    console.log(invoice_id)
    axios
      .get(`http://localhost:3128/invoices/user/${this.user.id}/${invoice_id}`, {
        headers: {
          "Token": token
        }
      })
      .then(res => {
        if (res.data.status == true) {
          this.transactions = res.data.transactions;
          this.invoice.id = res.data.invoice_id;
          let total = 0;
          this.transactions.forEach(element => {
            console.log(element)
            total += parseInt(element.price);
          });
          this.total_price = total;
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
  color: #426cb9;
}
.single-page {
  background-color: #ffffffe5;
}
.invoice {
  margin-top: 20px;
}
</style>
