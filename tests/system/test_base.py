from rabbitmqbeat import BaseTest


class Test(BaseTest):
    def test_base(self):
        """
        Basic test with exiting Rabbitmqbeat normally
        """
        self.render_config_template(modules=[{
            "name": "rabbitmq",
            "metricsets": ["overview"],
            "hosts": ["http://localhost:15672"],
            "username": "guest",
            "password": "guest",
            "period": "5s"
        }]
        )

        rabbitmqbeat_proc = self.start_beat()
        self.wait_until(lambda: self.output_lines() > 0)
        exit_code = rabbitmqbeat_proc.kill_and_wait()
        assert exit_code == 0
